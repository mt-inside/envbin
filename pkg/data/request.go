package data

import "net/http"

func getRequestData(r *http.Request) map[string]string {
    data := map[string]string{}

    data["RequestIP"] = r.RemoteAddr // This will be the last proxy; look at x-forwarded-for if you want to be better
    data["UserAgent"] = r.UserAgent()

    return data
}
