package main

import (
    "github.com/dangusev/goparser/app/utils"
    "github.com/dangusev/goparser/app"
    "log"
)

func main() {
    // Tune the logger
    log.SetFlags(log.LstdFlags | log.Lshortfile)

    context := &utils.GlobalContext{}
    context.PrepareSettings()

    app.RunCrawler(context)
    //
    //// Serve static
    //fs := http.FileServer(http.Dir(context.StaticDir))
    //http.Handle("/static/", http.StripPrefix("/static/", fs))
    //
    //// Routes
    //router := mux.NewRouter()
    //context.Router = router
    //router.Handle("/", utils.ExtendedHandler{GlobalContext: context, Get: handlers.MainHandler}).Name("main")
    //router.Handle("/api/queries/",
    //    utils.ExtendedHandler{GlobalContext: context,
    //        Get: handlers.QueriesListHandler,
    //        Post: handlers.QueriesAddHandler}).Name("queries-list")
    //router.Handle("/api/queries/{id}/",
    //    utils.ExtendedHandler{GlobalContext: context,
    //        Get: handlers.QueriesDetailHandler,
    //        Post: handlers.QueriesUpdateHandler,
    //        Delete: handlers.QueriesDeleteHandler}).Name("queries-detail")
    //router.Handle("/api/queries/{id}/items/",
    //    utils.ExtendedHandler{GlobalContext: context,
    //        Get: handlers.ItemsListHandler}).Name("items-list")
    //
    //http.Handle("/", router)
    //port := "8080"
    //log.Println(fmt.Sprintf("Run server on localhost:%s", port))
    //err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
    //if err != nil {
    //    log.Fatal(err)
    //}
}
