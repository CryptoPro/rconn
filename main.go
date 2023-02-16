package main

import (
	"net"
	"os"
)

func main() {
	initLogger()

	if len(os.Args) < 2 || (os.Args[1] != "-s" && os.Args[1] != "-c") {
		Log.Infof("Usage %s: [-s remote_port local_port | -c remote_addr remote_port local_addr local_port]", os.Args[0])
		os.Exit(1)
	}

	if os.Args[1] == "-s" {
		remoteListen, err := net.Listen("tcp", "0.0.0.0:"+os.Args[2])
		if err != nil {
			Log.Errorf("Error listening on remote port: %s", err.Error())
			os.Exit(1)
		}
		defer remoteListen.Close()

		localListen, err := net.Listen("tcp", "0.0.0.0:"+os.Args[3])
		if err != nil {
			Log.Errorf("Error listening on local port: %s", err.Error())
			os.Exit(1)
		}
		defer localListen.Close()

		statusRemote := make(chan bool)
		statusLocal := make(chan bool)

		for {
			remoteConn, err := remoteListen.Accept()
			if err != nil {
				Log.Errorf("Error connecting on remote port: %s", err.Error())
				os.Exit(1)
			}
			localConn, err := localListen.Accept()
			if err != nil {
				Log.Errorf("Error connecting on local port: %s", err.Error())
				os.Exit(1)
			}
			go pipeSocket(remoteConn, localConn, statusRemote)
			go pipeSocket(localConn, remoteConn, statusLocal)
		}

		// for {
		// 	status := <-statusLocal
		// 	if !status {
		// 		localConn, err = localListen.Accept()
		// 		if err != nil {
		// 			Log.Errorf("Error connecting on local port: %s", err.Error())
		// 			os.Exit(1)
		// 		}
		// 	}
		// 	go pipeSocket(false, statusLocal)
		// }
	}

	if os.Args[1] == "-c" {
		statusRemote := make(chan bool)
		//statusLocal := make(chan bool)

		for {
			remoteConn, err := net.Dial("tcp", os.Args[2]+":"+os.Args[3])
			if err != nil {
				Log.Errorf("Error dialing to remote port: %s", err.Error())
				continue
				//os.Exit(1)
			}

			go pipeSocket(remoteConn, nil, statusRemote)
			//go pipeSocket(localConn, remoteConn, statusLocal)
		}
		// for {
		// 	status := <-statusLocal
		// 	if !status {
		// 		localConn, err = net.Dial("tcp", os.Args[4]+":"+os.Args[5])
		// 		if err != nil {
		// 			Log.Errorf("Error dialing to local port: %s", err.Error())
		// 			os.Exit(1)
		// 		}
		// 	}
		// 	go pipeSocket(false, statusLocal)
		// }
	}
}

func pipeSocket(readConn net.Conn, writeConn net.Conn, status chan<- bool) {
	Log.Info("Started")
	for {
		//buf := make([]byte, 0x10000)
		buf := make([]byte, 0x1000)
		var err error
		var read, write int

		defer readConn.Close()

		read, err = readConn.Read(buf)

		if err != nil {
			Log.Errorf("Read error: %s", err.Error())
			status <- false
			return
		}
		Log.Info("Read: ", read)

		if writeConn == nil {
			writeConn, err = net.Dial("tcp", os.Args[4]+":"+os.Args[5])
			if err != nil {
				Log.Errorf("Error dialing to local port: %s", err.Error())
				return
				//os.Exit(1)
			}
			go pipeSocket(writeConn, readConn, status)
		}

		defer writeConn.Close()

		var start, written int
		for written < read {
			write, err = writeConn.Write(buf[start:read])
			if err != nil {
				Log.Errorf("Write error: %s", err.Error())
				status <- false
				return
			}
			Log.Info("Write: ", write)

			written += write
			start += write
		}
	}
	Log.Info("Finished")
}
