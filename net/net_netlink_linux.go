//go:build linux
// +build linux

package net

import (
	"context"
	"fmt"
	"strconv"
	"syscall"

	"github.com/shirou/gopsutil/v3/internal/common"
	"github.com/shirou/gopsutil/v3/net/netlink"
)

func connectionsPidMaxWithoutUidsWithContextByNetlink(_ context.Context, kind string, pid int32, max int, _ bool) ([]ConnectionStat, error) {
	tmap, ok := netConnectionKindMap[kind]
	if !ok {
		return nil, fmt.Errorf("invalid kind, %s", kind)
	}
	root := common.HostProc()
	var err error
	var inodes map[string][]inodeMap
	if pid == 0 {
		inodes, err = getProcInodesAll(root, max)
	} else {
		inodes, err = getProcInodes(root, pid, max)
	}
	if len(inodes) == 0 {
		// no inode for the pid
		return []ConnectionStat{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("cound not get pid(s), %d: %s", pid, err)
	}

	return connectionsNetLink(tmap, inodes)
}

func connectionsNetLink(kinds []netConnectionKindType, inodes map[string][]inodeMap) ([]ConnectionStat, error) {
	reply := make([][]ConnectionStat, len(kinds))
	errSlice := make([]error, len(kinds))
	var retErrr error

	for i, kind := range kinds {
		if kind.family == syscall.AF_UNIX {
			getUnixConnections(i, reply, errSlice)
		} else {
			getInetConnections(i, kind, reply, errSlice)
		}
	}

	for _, err := range errSlice {
		if err != nil {
			return nil, err
		}
	}

	ret := make([]ConnectionStat, 0)
	for i := range kinds {
		for _, cs := range reply[i] {
			iss, exists := inodes[strconv.Itoa(int(cs.Fd))]
			if exists {
				for _, is := range iss {
					cs.Pid = is.pid
					cs.Fd = is.fd
					ret = append(ret, cs)
				}
			}
		}
	}
	return ret, retErrr
}

func getInetConnections(i int, k netConnectionKindType, reply [][]ConnectionStat, errSlice []error) {
	msgs, err := netlink.InetConnections(uint8(k.family), k.netlinkProto)
	if err != nil {
		errSlice[i] = err
		return
	}
	t := make([]ConnectionStat, len(msgs))
	for i, msg := range msgs {
		conn := ConnectionStat{
			Family: uint32(msg.Family),
			Type:   uint32(k.sockType),
			Laddr:  Addr{msg.SrcIP().String(), uint32(msg.SrcPort())},
			Raddr:  Addr{msg.DstIP().String(), uint32(msg.DstPort())},
			Status: netlink.TCPState(msg.State).String(),
			Uids:   []int32{int32(msg.UID)},
			// Pid:    msg.Pid,
			Fd: msg.Inode,
		}
		t[i] = conn
	}
	reply[i] = t
}

func getUnixConnections(i int, reply [][]ConnectionStat, errSlice []error) {
	msgs, err := netlink.UnixConnections()
	if err != nil {
		errSlice[i] = err
		return
	}
	t := make([]ConnectionStat, len(msgs))
	for i, msg := range msgs {
		conn := ConnectionStat{
			Family: uint32(msg.Family),
			Status: netlink.TCPState(msg.State).String(),
			Uids:   []int32{},
			Laddr: Addr{
				IP: msg.Path,
			},
			Fd: msg.Inode,
		}
		t[i] = conn
	}
	reply[i] = t
}
