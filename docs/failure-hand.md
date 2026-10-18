1.in-memory failure handling

                 logQueue
                     ↓
                   worker
                     ↓
                PostgreSQL
                  ↙     ↘
              success   failure
                         ↓
                       retry
                         ↓
                     3 failures
                         ↓
                    failedQueue


2.