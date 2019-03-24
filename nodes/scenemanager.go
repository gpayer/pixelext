package nodes

var sceneManager *SceneManagerStruct

func SceneManager() *SceneManagerStruct {
	return sceneManager
}

type SceneManagerStruct struct {
	root nodeInternal
}

func (s *SceneManagerStruct) SetRoot(root nodeInternal) {
	if s.root != nil {
		s.root._unmount()
	}
	root._init()
	s.root = root
	root._mount()
}

func init() {
	sceneManager = &SceneManagerStruct{}
}
