# Go Concurrency Notes — LogFlow

## 1. The Most Important Mental Model: Go Code Is Sequential

Inside a single goroutine, ordinary Go statements execute **in order**.

```go
line1()
line2()
line3()
```

The program does:

```text
line1
  ↓
line2
  ↓
line3
```

It does **not** start `line3` while `line2` is still waiting.

For example:

```go
<-signalChan
close(shutdown)
```

The first line means:

> Wait until a value is received from `signalChan`.

Therefore:

```text
<-signalChan
    ↓
   WAIT
    ↓
signal arrives
    ↓
receive completes
    ↓
close(shutdown)
```

If the signal never arrives, `close(shutdown)` is never reached.

This is normal sequential execution.

---

# 2. Then Why Do We Need `select`?

Because sometimes we **don't want to wait for one specific thing**.

We want to wait for **multiple channel events**.

For example, our worker needs to deal with two possibilities:

```text
1. A log arrives in logQueue
2. Shutdown is requested
```

If we write:

```go
log := <-logQueue
```

the worker waits specifically for `logQueue`.

It cannot simultaneously wait for shutdown at that point.

That's where `select` comes in.

```go
select {
case log := <-logQueue:
    // process log

case <-shutdown:
    // handle shutdown
}
```

`select` means:

> **Wait for multiple channel operations and execute a case whose channel operation is ready.**

Conceptually:

```text
              select
             /      \
            ↓        ↓
       logQueue    shutdown
        ready?      ready?
```

The worker doesn't have to choose one channel beforehand.

It waits for either event.

---

# 3. `select` Does NOT Mean Everything Runs Simultaneously

This is important.

A `select` does **not** execute all cases at the same time.

Instead:

```text
select
  ↓
check which channel operations are ready
  ↓
choose one ready case
  ↓
execute that case
```

If multiple cases are ready, Go may choose any one of the ready cases.

So don't think:

> "Both cases execute simultaneously."

Think:

> **"I'm waiting on multiple possible events."**

---

# 4. Why Couldn't We Just Use `if/else`?

Consider:

```go
if logQueue has a log {
    process it
} else {
    shutdown
}
```

The problem is that `if/else` is ordinary sequential logic.

The worker would have to check the queue at one particular moment.

But what if:

```text
queue isn't ready
shutdown isn't ready
```

and the worker needs to **wait**?

Channels and `select` are designed for this kind of synchronization.

With:

```go
select {
case log := <-logQueue:
    // ...
case <-shutdown:
    // ...
}
```

the worker can wait until one of those channel operations becomes possible.

---

# 5. Our `signalChan`

Now let's understand the three shutdown lines properly.

```go
signalChan := make(chan os.Signal, 1)

signal.Notify(signalChan, os.Interrupt)

<-signalChan
```

## Line 1

```go
signalChan := make(chan os.Signal, 1)
```

### `signalChan`

A variable that refers to a channel.

### `chan os.Signal`

The channel carries values of type `os.Signal`.

### `1`

The channel has a buffer capacity of one signal.

So we can imagine:

```text
signalChan
┌─────────────────┐
│  1 signal slot  │
└─────────────────┘
```

We're creating a small mailbox for OS signals.

---

# 6. What Is `signal` in `signal.Notify`?

This can be confusing:

```go
signal.Notify(signalChan, os.Interrupt)
```

`signal` is the **`os/signal` package**.

`Notify` is a **function from that package**.

So:

```text
signal
   ↓
package

Notify
   ↓
function

signalChan
   ↓
our channel

os.Interrupt
   ↓
the OS signal we care about
```

The line means roughly:

> Tell Go to send an `os.Interrupt` signal into `signalChan` when the operating system generates one.

For us, that normally happens when we press:

```text
Ctrl + C
```

---

# 7. What Does `<-signalChan` Do?

```go
<-signalChan
```

The `<-` here means **receive from the channel**.

We are not storing the received value because we don't care about the actual signal value.

If we wanted to store it:

```go
sig := <-signalChan
```

Then `sig` would contain the received signal.

But we only need to know:

> **"Did an interrupt arrive?"**

So:

```go
<-signalChan
```

means:

> **Wait here until a signal arrives in `signalChan`.**

---

# 8. The Three Lines Together

Now put everything together:

```go
signalChan := make(chan os.Signal, 1)
signal.Notify(signalChan, os.Interrupt)
<-signalChan
```

The flow is:

```text
Create channel
     ↓
Tell Go which OS signal to put there
     ↓
main waits for something to arrive
     ↓
        Ctrl + C
           ↓
      OS interrupt
           ↓
     signalChan
           ↓
    <-signalChan
           ↓
   receive completes
           ↓
   next line executes
```

That last part is important.

Because ordinary code is sequential:

```go
<-signalChan
close(shutdown)
```

means:

```text
WAIT
 ↓
signal arrives
 ↓
receive completes
 ↓
continue to next statement
 ↓
close(shutdown)
```

If no signal arrives:

```text
<-signalChan
    ↓
  WAIT...
    ↓
  WAIT...
    ↓
close(shutdown) ❌
```

`close(shutdown)` is never executed.

---

# 9. Why `close(shutdown)`?

We created:

```go
var shutdown = make(chan struct{})
```

Its type is:

```text
chan struct{}
```

`struct{}` is an empty struct.

It contains no fields and therefore carries no useful data.

We don't need to send information.

We only need a signal that means:

> **"Shutdown has been requested."**

So instead of:

```go
shutdown <- something
```

we use:

```go
close(shutdown)
```

Closing the channel becomes our broadcast signal.

---

# 10. Why Is Closing Useful With Multiple Workers?

We have:

```text
Worker 1
Worker 2
Worker 3
```

All of them can listen to:

```go
<-shutdown
```

When `main` does:

```go
close(shutdown)
```

the closed state can be observed by all workers.

Conceptually:

```text
              close(shutdown)
                    ↓
          ┌─────────┼─────────┐
          ↓         ↓         ↓
       Worker 1  Worker 2  Worker 3
          ↓         ↓         ↓
       shutdown  shutdown  shutdown
       detected  detected  detected
```

We aren't sending three separate messages.

We're changing the state of the shared channel to **closed**.

---

# 11. The Complete Shutdown Chain

Our LogFlow shutdown mechanism is therefore:

```text
                 Ctrl + C
                    ↓
             Operating System
                    ↓
              os.Interrupt
                    ↓
             signal.Notify()
                    ↓
               signalChan
                    ↓
              <-signalChan
                    ↓
             main continues
                    ↓
             close(shutdown)
                    ↓
          ┌─────────┼─────────┐
          ↓         ↓         ↓
       Worker 1  Worker 2  Worker 3
          ↓         ↓         ↓
       finish    finish    finish
       current   current   current
          ↓         ↓         ↓
              drain logQueue
                    ↓
              queue becomes empty
                    ↓
              workers exit
                    ↓
              WaitGroup done
                    ↓
              pool.Close()
                    ↓
                 program exits
```

---

# 12. The Big Picture

Remember these four concepts together:

### Sequential execution

```text
A → B → C
```

One goroutine normally executes statements in order.

### Goroutines

```text
A ──→ goroutine 1
B ──→ goroutine 2
C ──→ goroutine 3
```

They allow independent work to execute concurrently.

### Channels

```text
goroutine 1 ──→ channel ──→ goroutine 2
```

They allow goroutines to communicate and synchronize.

### `select`

```text
             select
            /      \
       channel A  channel B
```

It lets one goroutine **wait on multiple channel operations**.

That is why `select` is important in our LogFlow worker.

We don't use `select` because "Go code needs select."

We use it because **our worker has multiple things it needs to react to**:

```text
"Did a log arrive?"
       OR
"Was shutdown requested?"
```

And that is the real reason behind the syntax.