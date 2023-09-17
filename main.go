package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"path"
	"time"
)

// Next:
// 	Make a gui for submitting stuff?
// 	New format for files?

// Todo:
//	Audit

// Decisions:
//   - there will be no hot reloading, not worth it, just restart the thing, nobody will notice.
//   - anti matching of tags with !tag has no real use at this point.
var (
	address         = flag.String("ip", "127.0.0.1", "the address to listen on")
	port            = flag.String("port", "8000", "the port the app listens on")
	entriesPerPage  = flag.Uint("epp", 20, "the amount of entries displayed per page") //fixme make higher for release
	applicationdata = flag.String("appdata", "user_applicationdata", "the directory containing all the app_*.csv files")
)

func main() {
	flag.Parse()
	if *entriesPerPage < 1 {
		panic("you need to have at least one entry per page")
	}

	apps, err := NewApplicationsFromPath(*applicationdata)
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()

	// App details resources
	screenshots := http.FileServer(http.Dir(path.Join(*applicationdata, "screenshots")))
	mux.Handle("/screenshots/", http.StripPrefix("/screenshots/", screenshots))

	// CSS and favicon for the html files
	mux.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("user_templates/resources"))))

	// Search API for the Frontend
	mux.Handle("/api/v1/search", searchhandler(apps))

	// Routing detail
	mux.HandleFunc("/search",
		searchedpage(apps,
			"user_templates/index.html"),
	)
	mux.HandleFunc("/app/",
		applicationparser(apps,
			"user_templates/apptemplate.html"),
	) // /app/{{appname}}
	//hint: since i convenience-redirected non existing urls on / to search,
	//you wont be able to "search" for an application called "about", but an
	//application called about works, because of the /app/ prefix.
	//would be stupid to cap the feature because of that and not enough gain
	//for the extra overhead for handling it inside "/".
	mux.HandleFunc("/about",
		aboutpage(apps,
			"user_templates/about.html"),
	)
	mux.HandleFunc("/",
		indexpage(apps,
			"user_templates/index.html", *entriesPerPage),
	)

	srv := &http.Server{
		Handler:      mux,
		Addr:         net.JoinHostPort(*address, *port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("Ready... http://", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
