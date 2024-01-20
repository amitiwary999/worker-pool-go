# worker-pool-go

My implementation to try worker pool. Initial version is:
 n number of goroutine is running and we keep adding the job in queue and any one goroutine take job from queue and do it.
 Buffered channel is helpful where there is always some task on buffer even if all goroutine is busy.