/logs/payment - tells us the path

/logs?service=payment - tells us the query path 


1./logs?service=payment&level=ERROR

gives us like  r.URL.Query() , so we need r.URL.Query().Get("service")