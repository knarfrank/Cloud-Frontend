package main


import "net/http"
import "net/url"
import "io"
import "fmt"
import "bytes"
import "time"
import "encoding/json"

import "google.golang.org/appengine"
import "google.golang.org/appengine/taskqueue"
import "google.golang.org/appengine/datastore"
import "google.golang.org/appengine/urlfetch"
import "golang.org/x/net/context"

type ScanData struct {
    Url          string
    Ip           string
    CommonName   string
    Subject      string
    ValidFrom    time.Time
    ValidTill    time.Time
    Fingerprint  string
    SignatureAlgorithm string
    PublicKeyAlgorithm string
    PublicE             int
}






func init() {
    fs := http.FileServer(http.Dir("frontend"))
    http.Handle("/", fs)
    http.HandleFunc("/queue", handleQueue)
    http.HandleFunc("/scan", handleScan)
    http.HandleFunc("/store", handleStore)
    http.HandleFunc("/get", handleGet)
    http.HandleFunc("/getrecent", handleGetAll)

}

func handleScan(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
    ctx := appengine.NewContext(r)
    client := urlfetch.Client(ctx)
    resp, err := client.Get("http://reconflex-154712.appspot-preview.com/scan?url="+r.FormValue("url"))
    if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
    }
    fmt.Fprintf(w, "HTTP GET returned status %v", resp.Status)
}

func handleQueue(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
    io.WriteString(w, "Attempting to queue scan!\n" + r.FormValue("url"))


    var ctx context.Context
    ctx = appengine.NewContext(r)
    t := taskqueue.NewPOSTTask("/scan", url.Values{
	        "url": {r.FormValue("url")},
	})
    if t == nil {
        io.WriteString(w, "CRAP\n")
        return;
    }
    if ctx == nil {
        io.WriteString(w, "POOP\n")
        return;
    }
	// Use the transaction's context when invoking taskqueue.Add.
	_, err := taskqueue.Add(ctx, t, "scan-queue")
	if err != nil {
		io.WriteString(w, err.Error())
	}
}


func handleGet(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
    ctx := appengine.NewContext(r)
    q := datastore.NewQuery("scan").Filter("Url =", r.FormValue("url"))
    var arrayData []ScanData

    b := new(bytes.Buffer)
    for t := q.Run(ctx); ; {
        var x ScanData
        _, err := t.Next(&x)
        if err == datastore.Done {
            break
        }
        if err != nil {
            io.WriteString(w, err.Error())
            return
        }
        arrayData = append(arrayData, x)

    }
    json, _ := json.Marshal(arrayData)
    fmt.Fprintf(b, "%s", json)
    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    io.Copy(w, b)
}


func handleGetAll(w http.ResponseWriter, r *http.Request) {

    ctx := appengine.NewContext(r)
    q := datastore.NewQuery("scan")
    var arrayData []ScanData

    b := new(bytes.Buffer)
    for t := q.Run(ctx); ; {
        var x ScanData
        _, err := t.Next(&x)
        if err == datastore.Done {
            break
        }
        if err != nil {
            io.WriteString(w, err.Error())
            return
        }
        arrayData = append(arrayData, x)

    }
    json, _ := json.Marshal(arrayData)
    fmt.Fprintf(b, "%s", json)
    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    io.Copy(w, b)
}



func handleStore(w http.ResponseWriter, req *http.Request) {
    var ctx context.Context
    ctx = appengine.NewContext(req)

    decoder := json.NewDecoder(req.Body)
    var t ScanData
    err := decoder.Decode(&t)
    if err != nil {
        io.WriteString(w, err.Error());
    }
    defer req.Body.Close()


    key := datastore.NewIncompleteKey(ctx, "scan", nil)
    if _, err := datastore.Put(ctx, key, &t); err != nil {
          io.WriteString(w, err.Error());
    }
}
