redo all schemas
- index
- review data types
- batch updates for pty sessions? (maybe every 5-10secs?) or at least wait till the first conn is done
- maybe denormalise pty session + pty conn if ppl don't really use rejoin + observe
- if CR synced to our own db, can use CR number as foreign key for pty session start conn cr so can preload
- maybe pty sessions need to be partitioned 
- 