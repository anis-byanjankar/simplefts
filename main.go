package main

import (
	"flag"
	"fmt"
	"log"
	"time"
	"net/http"
	"net/http/httputil"
)

func main() {
	var dumpPath, query string
	flag.StringVar(&dumpPath, "p", "enwiki-latest-abstract1.xml.gz", "wiki abstract dump path")
	flag.StringVar(&query, "q", "Small wild cat", "search query")
	flag.Parse()

	log.Println("Starting simplefts")

	start := time.Now()
	docs, err := loadDocuments(dumpPath)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Loaded %d documents in %v", len(docs), time.Since(start))

	start = time.Now()
	idx := make(index)
	idx.add(docs)
	log.Printf("Indexed %d documents in %v", len(docs), time.Since(start))

	start = time.Now()
	matchedIDs := idx.search(query)
	log.Printf("Search found %d documents in %v", len(matchedIDs), time.Since(start))

	for _, id := range matchedIDs {
		doc := docs[id]
		log.Printf("%d\t%s\n", id, doc.Text)
	}
	// handle route using handler function
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        	fmt.Fprintf(w, "Welcome to new server!")
		
		reqDump, err := httputil.DumpRequest(r, true)
		if err != nil {
			log.Fatal(err)
		}
		searchQuery := r.URL.Path[1:]

		searchedIDs := idx.search(searchQuery)
		res := ""
		for _, idN := range searchedIDs {
			doc := docs[idN]
			res += fmt.Sprintf("%d\t%s\n", idN, doc.Text) + "\n\n\n\n"
			//res += string(int(idN)) + " " + doc.Text+ "\n\n\n\n"
		}

		fmt.Fprintf(w,"REQUEST:\n%s", string(reqDump))
		fmt.Fprintf(w,"REQUEST:\n%s", r.URL.Path[1:])
		fmt.Fprintf(w,"Search Result:\n%s",res)


	})

        // listen to port
    	http.ListenAndServe(":5050", nil)

}
