package main

import "net"

type rooms struct{

	name string 
	member map[net.Addr]*client
	len int 
}

func ( r *rooms) broadcast(sender *client , msg string){

	for addr,m :=range r.member{

		if sender.conn.RemoteAddr() != addr {
			m.msg(msg)
		}
	}
}