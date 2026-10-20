the reason we have select is that , 

without select if we use if & else : prolly , when both http handlers goroutines gets in and sees it , they be like
seeing the same , 

for example :

But now imagine two HTTP requests arrive at almost the same time:
Request A                  Request B
    │                          │
    ▼                          ▼
"Is queue full?"           "Is queue full?"
    │                          │
    NO                         NO
    │                          │
    ▼                          ▼
"Okay, I'll send"          "Okay, I'll send"

Both checked before either one actually sent.

so , inroder to avoid this , we don't need to check , we need to DIRECTLY INSERT ITTT

that's why select 

select {
case logQueue <- incomingLog:
    // SUCCESS
default:
    // COULDN'T SEND IMMEDIATELY
} 

whyy DEFAULT :
 because without default the http handler goroutine be waiting on it , which is worse because nah , bruh , it doesn't have to bee
 so , we'll send our default : so the handler doesn't have to waitt

