package logger

// SubContext is a sub-reference to a logger with a specific heading and set of default labels for messages.
// It implements the full logger suite but forwards them up to the parent logger.
type SubContext struct {
	log         *Logger
	headings    []string
	labels      map[string]string
	annotations map[string]string
}

// Logger returns the underlying logger.
func (sc *SubContext) Logger() *Logger {
	return sc.log
}

// SubContext returns a further sub-context with a given heading.
func (sc *SubContext) SubContext(heading string) *SubContext {
	return &SubContext{
		headings:    append(sc.headings, heading),
		labels:      sc.labels,
		annotations: sc.annotations,
	}
}

// Headings returns the headings.
func (sc *SubContext) Headings() []string {
	return sc.headings
}

// WithLabel adds a label.
func (sc *SubContext) WithLabel(key, value string) *SubContext {
	if sc.labels == nil {
		sc.labels = map[string]string{}
	}
	sc.labels[key] = value
	return sc
}

// WithLabels sets the labels.
func (sc *SubContext) WithLabels(labels map[string]string) *SubContext {
	sc.labels = labels
	return sc
}

// Labels returns the sub-context labels.
func (sc *SubContext) Labels() map[string]string {
	return sc.labels
}

// WithAnnotations sets the annotations.
func (sc *SubContext) WithAnnotations(annotations map[string]string) *SubContext {
	sc.annotations = annotations
	return sc
}

// WithAnnotation adds an annotation.
func (sc *SubContext) WithAnnotation(key, value string) *SubContext {
	if sc.annotations == nil {
		sc.annotations = map[string]string{}
	}
	sc.annotations[key] = value
	return sc
}

// Annotations returns the sub-context annotations.
func (sc *SubContext) Annotations() map[string]string {
	return sc.annotations
}

// Sillyf writes a message.
func (sc *SubContext) Sillyf(format string, args ...Any) {
	msg := Messagef(Silly, format, args...).WithHeadings(sc.headings...)
	sc.injectLabels(msg)
	sc.injectAnnotations(msg)
	sc.log.Trigger(msg)
}

// SyncSillyf synchronously writes a message.
func (sc *SubContext) SyncSillyf(format string, args ...Any) {
	msg := Messagef(Silly, format, args...).WithHeadings(sc.headings...)
	sc.injectLabels(msg)
	sc.injectAnnotations(msg)
	sc.log.SyncTrigger(msg)
}

// Infof writes a message.
func (sc *SubContext) Infof(format string, args ...Any) {
	msg := Messagef(Info, format, args...).WithHeadings(sc.headings...)
	sc.injectLabels(msg)
	sc.injectAnnotations(msg)
	sc.log.Trigger(msg)
}

// SyncInfof synchronously writes a message.
func (sc *SubContext) SyncInfof(format string, args ...Any) {
	msg := Messagef(Info, format, args...).WithHeadings(sc.headings...)
	sc.injectLabels(msg)
	sc.injectAnnotations(msg)
	sc.log.SyncTrigger(msg)
}

// Debugf writes a message.
func (sc *SubContext) Debugf(format string, args ...Any) {
	msg := Messagef(Debug, format, args...).WithHeadings(sc.headings...)
	sc.injectLabels(msg)
	sc.injectAnnotations(msg)
	sc.log.Trigger(msg)
}

// SyncDebugf synchronously writes a message.
func (sc *SubContext) SyncDebugf(format string, args ...Any) {
	msg := Messagef(Debug, format, args...).WithHeadings(sc.headings...)
	sc.injectLabels(msg)
	sc.injectAnnotations(msg)
	sc.log.SyncTrigger(msg)
}

// Warningf writes a message.
func (sc *SubContext) Warningf(format string, args ...Any) {
	msg := Errorf(Warning, format, args...).WithHeadings(sc.headings...)
	sc.injectLabels(msg)
	sc.injectAnnotations(msg)
	sc.log.Trigger(msg)
}

// SyncWarningf synchronously writes an error message.
func (sc *SubContext) SyncWarningf(format string, args ...Any) {
	msg := Errorf(Warning, format, args...).WithHeadings(sc.headings...)
	sc.injectLabels(msg)
	sc.injectAnnotations(msg)
	sc.log.SyncTrigger(msg)
}

// Errorf writes an error  message.
func (sc *SubContext) Errorf(format string, args ...Any) {
	msg := Errorf(Error, format, args...).WithHeadings(sc.headings...)
	sc.injectLabels(msg)
	sc.injectAnnotations(msg)
	sc.log.Trigger(msg)
}

// SyncErrorf synchronously writes an error message.
func (sc *SubContext) SyncErrorf(format string, args ...Any) {
	msg := Errorf(Error, format, args...).WithHeadings(sc.headings...)
	sc.injectLabels(msg)
	sc.injectAnnotations(msg)
	sc.log.SyncTrigger(msg)
}

// Fatalf writes an error  message.
func (sc *SubContext) Fatalf(format string, args ...Any) {
	msg := Errorf(Fatal, format, args...).WithHeadings(sc.headings...)
	sc.injectLabels(msg)
	sc.injectAnnotations(msg)
	sc.log.Trigger(msg)
}

// SyncFatalf synchronously writes an error message.
func (sc *SubContext) SyncFatalf(format string, args ...Any) {
	msg := Errorf(Fatal, format, args...).WithHeadings(sc.headings...)
	sc.injectLabels(msg)
	sc.injectAnnotations(msg)
	sc.log.SyncTrigger(msg)
}

func (sc *SubContext) injectLabels(e Event) {
	if sc.labels == nil {
		return
	}
	if typed, isTyped := e.(EventLabels); isTyped {
		for key, value := range sc.labels {
			typed.Labels()[key] = value
		}
	}
}

func (sc *SubContext) injectAnnotations(e Event) {
	if sc.annotations == nil {
		return
	}
	if typed, isTyped := e.(EventAnnotations); isTyped {
		for key, value := range sc.annotations {
			typed.Annotations()[key] = value
		}
	}
}
