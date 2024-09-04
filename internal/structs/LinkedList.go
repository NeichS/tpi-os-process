package structs

type Node struct {
    Process *Process
    Next    *Node
}

type LinkedList struct {
    Head *Node
}

func NewLinkedList() *LinkedList {
	return &LinkedList{}
}

func (l *LinkedList) Append(p *Process) {
    newNode := &Node{Process: p}
    if l.Head == nil {
        l.Head = newNode
        return
    }
    current := l.Head
    for current.Next != nil {
        current = current.Next
    }
    current.Next = newNode
}

func (l *LinkedList) Remove(name string) {
    if l.Head == nil {
        return
    }
    if l.Head.Process.Name == name {
        l.Head = l.Head.Next
        return
    }
    current := l.Head
    for current.Next != nil {
        if current.Next.Process.Name == name {
            current.Next = current.Next.Next
            return
        }
        current = current.Next
    }
}

func (l *LinkedList) Find(name string) *Process {
    current := l.Head
    for current != nil {
        if current.Process.Name == name {
            return current.Process
        }
        current = current.Next
    }
    return nil
}