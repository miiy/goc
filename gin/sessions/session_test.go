package sessions

import ginsessions "github.com/gin-contrib/sessions"

type testSession struct {
	values  map[interface{}]interface{}
	flashes map[string][]interface{}
	saveErr error
	saved   bool
}

func newTestSession() *testSession {
	return &testSession{
		values:  make(map[interface{}]interface{}),
		flashes: make(map[string][]interface{}),
	}
}

func (s *testSession) ID() string { return "" }

func (s *testSession) Get(key interface{}) interface{} {
	return s.values[key]
}

func (s *testSession) Set(key interface{}, val interface{}) {
	s.values[key] = val
}

func (s *testSession) Delete(key interface{}) {
	delete(s.values, key)
}

func (s *testSession) Clear() {
	for key := range s.values {
		delete(s.values, key)
	}
}

func (s *testSession) AddFlash(value interface{}, vars ...string) {
	key := "_flash"
	if len(vars) > 0 {
		key = vars[0]
	}
	s.flashes[key] = append(s.flashes[key], value)
}

func (s *testSession) Flashes(vars ...string) []interface{} {
	key := "_flash"
	if len(vars) > 0 {
		key = vars[0]
	}

	values := s.flashes[key]
	delete(s.flashes, key)
	return values
}

func (s *testSession) Options(ginsessions.Options) {}

func (s *testSession) Save() error {
	s.saved = true
	return s.saveErr
}
