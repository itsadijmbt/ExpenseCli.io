                      +--------------------------------+
                      |      Bubble Tea Runtime       |
                      |  (p.eventLoop, see below)     |
                      +--------------------------------+
                                   ^
                                   |  (6) Cmds from Update loop
                                   |
                              +----+----+
                              |  Cmd    |
                              | func()  |
                              | tea.Msg |
                              +----+----+
                                   |
                                   v
(1) User presses a key or the window resizes
       │
       ▼
+-------------+    (2) Runtime wraps it into      +-------------------+
| Terminal    |───> tea.KeyMsg / tea.WindowSizeMsg│                   │
+-------------+    and sends it on p.msgs channel +-------------------+
                                                      │
                                                      ▼
                                              +---------------+
                                              | model.Update  |
                                              | (msg)         |
                                              +---------------+
                                                      │
                                                      │ returns
                                                      │ (newModel, cmd)
                                                      ▼
                                              +---------------+
                                              | model.View    |
                                              | (render output)|
                                              +---------------+
                                                      │
                                                      │
                                                      ▼
                                            Displays updated UI
                                                      │
                                                      └─────┐
                                                            │
                         (5) bubble’s Cmd is sent to cmds ──┘
                         channel and executed asynchronously :contentReference[oaicite:0]{index=0}
                                                            │
                                                            ▼
               +-------------------------------+   (7) when it finishes it
               |     tea.Cmd runs in a         |   returns a tea.Msg (e.g.  |
               |     goroutine:                |──> newExpenseMsg, TickMsg, |
               |     fn() tea.Msg              |    errMsg, etc.             |
               +-------------------------------+
                                                            │
                                                            ▼
                                                 Back to step (2):
                                          runtime receives that Msg too


////////////////////////////////INFO//////////////////////////////
tea.Cmd lets you kick off work in the background (timers, I/O, calculations, etc.) without blocking the main UI loop.
Bubble Tea queues that function,runs it in its own goroutine, and when it finishes it sends its returned tea.Msg back into your Update method
That way you maintain a steady pipeline of messages driving your UI updates, all without freezing the interface.


/////////////////////////TODO//////////////////////////
Read Code First: Running and stepping through the official examples is worth more than any blog post.

Small Iterations: Start with a tiny app—just a text input or just a spinner—then grow it.

Leverage Type Safety: If you ever see a tea.Model in your code, remember you’ll often need a type‑assertion back to your concrete bubble.

Stick with it one small feature at a time, and soon this will feel as natural as any web‑framework MVU. Good luck!

You said:
. Deeper Dives & Patterns
Composability Patterns
Study the “composable‑views” example in the Bubble Tea repo to see how multiple bubbles (lists, spinners, paginators, help menus) fit together under one root model.

State Management Techniques
Look at how flags or small enums let you switch between “pages” (the pattern you need for a prompt sub‑view).
link for these


ChatGPT said:

Thought for a couple of seconds
Here are the direct links:

Composable Views Example
See how multiple bubbles are composed under one root model:
https://github.com/charmbracelet/bubbletea/blob/master/examples/composable-views/main.go 
GitHub

State Management Pattern (Enum + Switch)
In the same example, sessionState is a small enum used to switch between “pages” (timer vs spinner):
https://github.com/charmbracelet/bubbletea/blob/master/examples/composable-views/main.go#L22-L31 
GitHub

These illustrate how you can structure your own app with multiple sub‑models and switch view modes via flags or enums.