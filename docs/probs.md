1.What happens if the worker later fails to insert the log?

and The client already received:
```
{"status":"accepted"} ,


 2.or letting WAITGROUP and once it is done : we'll let it know we're done and 
 
 
 3.Do we want the HTTP request to wait when the queue is full, or reject/drop the log?
```

4.worker processing(post) a log whihc is being read by get request . so , what happends, consistency or dirty read?


| Problem | Our design question |
|---|---|
| Worker DB failure | How do we mark/retry `FAILED`? |
| POST says accepted too early | What exactly does `accepted` mean? |
| Queue full | Backpressure / reject / drop? |
| GET during INSERT | PostgreSQL prevents dirty reads; visibility depends on commit |
| Need to know persistence | Track processing state / query by ID |

5.Should LogFlow refuse to start when PostgreSQL is unavailable, or should it start and queue logs until the database comes back?