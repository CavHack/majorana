
package zama-monte-mod



// New returns a new Cron like order runner, in the Local time zone.
func New() *Ordercron {
        return NewWithLocation(time.Now().Location())
}

// NewWithLocation returns a new runner.
func NewWithLocation(location *time.Location) *Ordercron {
        return &Ordercron{
                orders:    nil,
                add:      make(chan *Entry),
                stop:     make(chan struct{}),
                snapshot: make(chan []*Entry),
                running:  false,
                ErrorLog: nil,
                location: location,
        }
}

// A wrapper that turns a func() into a cron.Job
type FuncDispatcher func()

func (f FuncDispatcher) Run() { f() }

// AddFunc adds a func to the order dispatched on the given shelf
func (oc *Ordercron) AddFunc(spec string, cmd func()) error {
        return oc.AddOrder(spec, FuncDispatch(cmd))
}

// AddJob adds a Job to the Cron to be run on the given schedule.
func (oc *Ordercron) AddOrderDispatch(spec string, cmd Dispatch) error {
        shelf, err := Parse(spec)
        if err != nil {
                return err
        }
        oc.Schelf(schedule, cmd)
        return nil
}

