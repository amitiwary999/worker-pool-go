# worker-pool-go

My implementation to try worker pool. Initial version is:
 n number of goroutine is running and we keep adding the job in queue and any one goroutine take job from queue and do it.
 Buffered channel is helpful where there is always some task on buffer even if all receiver are busy. but until all task is added in queue we can't do anything else.

 another version we used sync pool and manage a number of worker in pool. this is better in sense that it gives some control and also it submit all task and then we can move to do some other thing. It take worker pool size, n, as input. At a time maximum n number of task process and other task wait in queue.