package main

func main() {

}

type Beacon struct {
	Config *jobutil.JobConfig
}

// Name returns the job name.
func (b *Beacon) Name() string {
	return config.Name
}
