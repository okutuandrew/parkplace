package Middlewares 

import  (
"net/http"
"github.com/gorilla/sessions"
"fmt"

)

var store = sessions.NewCookieStore([]byte("super-secret-key")) // 16/24/32 bytes

func  SessionTracker(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		  session, err := store.Get(r, "driver-session")
        if err != nil {
            // Optional: log error, but still continue or block
            fmt.Println("session error:", err)
        }

        // Example: log whether user is authenticated
        if auth, ok := session.Values["authenticated"].(bool); ok && auth {
            fmt.Println("authenticated driver:", session.Values["driverID"])
        } else {
            fmt.Println("anonymous or not authenticated")
        }

		
    next.ServeHTTP(w, r)	
		
	})
}