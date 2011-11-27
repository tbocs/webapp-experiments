package tbocsincloud 

import ("template"
        "http"
        "appengine"
        //"appengine/datastore"
        "appengine/channel"
        "rand"
        "strconv"
)

var num_clients int = 0
var map_clients = make(map[string]string)

func init () {
  http.HandleFunc("/", root)
  http.HandleFunc("/callback", client_callback)
  http.HandleFunc("/_ah/channel/connected/", sys_channel_connected)
  http.HandleFunc("/_ah/channel/disconnected/", sys_channel_disconnected)
}

func root (w http.ResponseWriter, r *http.Request) {
  c := appengine.NewContext(r)
  clientTemplate := template.Must(template.ParseFile("client.html"))
  token_value := strconv.Itoa(rand.Int())
  token_key, _ := channel.Create(c, token_value)
  clientTemplate.Execute(w, token_key)
  map_clients[token_key] = token_value;
}

type MSG struct {
  Status bool
  Nickname string
  Message string
}

func client_callback (w http.ResponseWriter, r *http.Request) {
  c := appengine.NewContext(r)
  itsStatus, _ := strconv.Atob(r.FormValue("Status"))
  msg := MSG{itsStatus, r.FormValue("Nickname"), r.FormValue("Message")}
  for _, value := range map_clients {
      channel.SendJSON(c, value, msg)
  }
}

func sys_channel_connected (w http.ResponseWriter, r *http.Request) {
  //c := appengine.NewContext(r)
  num_clients ++
}

func sys_channel_disconnected (w http.ResponseWriter, r *http.Request) {
  //c := appengine.NewContext(r)
  num_clients --
  if (num_clients == 0) {
    map_clients = make(map[string]string)
  }
  
 
}
