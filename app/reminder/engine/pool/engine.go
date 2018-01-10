package pool

/*
	get all time intervals
	get all user data
	create pull of goroitine workers
	sort intervals
	put them into stack
	pick one and create timer
	listen for a timer
	when arised put all jobs in a queue
	then create next timer of (current duration - prev duration)

	worker get a job from a queue
	post request and check response
	if response have a "paid" flag worker write to "paid" channel the user id (?)
	goroutine listen for a paid users
	it receives user id and set flag to the user in a user data.
*/

type Engine struct {
	host string
}

func (e *Engine) Init(data [][]string, host string) error {
	e.host = host
	return nil
}

func (e Engine) Process() error {
	return nil
}
