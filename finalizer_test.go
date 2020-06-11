package txmanager

import "testing"

func TestCommit(t *testing.T) {
	var finalized, committed bool
	finalizer := MakeFinalizer("test")
	finalizer.Register(
		func() error { finalized = true; return nil },
		func() { committed = true },
		nil,
	)
	err := finalizer.Finalize()
	if err != nil {
		t.Fatalf(err.Error())
	}
	if !finalized {
		t.Fatalf("finalizer not executed")
	}
	if committed {
		t.Fatalf("commit function incorreclty executed")
	}
	finalizer.Commit()
	if !committed {
		t.Fatalf("commit function no executed")
	}
}
