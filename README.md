# worker-pool-go

My implementation to try worker pool. Initial version is:
 n number of goroutine is running and we keep adding the job in queue and any one goroutine take job from queue and do it.
 Buffered channel is helpful where there is always some task on buffer even if all receiver are busy. but until all task is added in queue we can't do anything else.

 another version:
  we used sync pool and manage a number of worker in pool. It take worker pool size, n, as input. At a time maximum n number of task process. When a worker is done, it is added again back to pool for other task. After expiry time we release the goroutine to save the resource. This is inspired from the `ants` library.