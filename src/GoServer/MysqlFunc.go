// MysqlFunc
package main

import (
	"protobuf"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/protobuf/proto"
)

type PlayerMysqlInfo struct {
	Userid    int32  //玩家唯一id
	Account   string //玩家的账号
	Passwd    string //玩家密码信息
	Name      string //玩家昵称
	Autograph string //玩家签名
	Grade     int32  //玩家段位
	Money     int32  //玩家钱币
	Icon      int32  //玩家头像id
	Sex       int32  //玩家性别
	Gnum      int32  //玩家总局数
	Gsnum     int32  //玩家升降局胜局
	Gfnum     int32  //玩家升降局输局
	Gdnum     int32  //玩家升降局平局的局数
	Fgsnum    int32  //玩家友谊局胜利的局数
	Fgfnum    int32  //玩家友谊局失败的局数
	Fgdnum    int32  //玩家友谊局平局的局数
}

func PlayerLogin(account string, passwd string) (*PlayerMysqlInfo, protobuf.ErrCodeType) {
	sql := "select p.id,p.passwd,p.name,p.autograph,p.grade,p.money,p.icon,p.sex,p.gnum,p.gsnum,p.gfnum,p.gdnum,p.fgsnum,p.fgfnum,p.fgdnum from yfgo.player_info p where p.account = ?"
	row := server.mysqlDB.QueryRow(sql)
	info := &PlayerMysqlInfo{}
	err := row.Scan(&info.Userid, &info.Passwd, &info.Name, &info.Autograph, &info.Grade,
		&info.Money, &info.Icon, &info.Sex, &info.Gnum, &info.Gsnum, &info.Gfnum,
		&info.Gdnum, &info.Fgsnum, &info.Fgfnum, &info.Fgdnum)
	log.Println(info)
	if err != nil {
		log.Error("mysql scan error:", err)
		return nil, protobuf.ErrCodeType_MysqlError
	}
	if info.Userid == 0 {
		return nil, protobuf.ErrCodeType_AccountNotExistError
	}
	if passwd != info.Passwd {
		return nil, protobuf.ErrCodeType_AccountOrPasswdNotMatchError
	}
	return info, protobuf.ErrCodeType_Success
}

func AddRegisterPlayerInfo(req *protobuf.PlayerRegisterRequest) (int32, protobuf.ErrCodeType) {
	sql := "select count(*) from yfgo.player_info p where p.account = ?"
	var (
		num    int32 = -1
		userid int32 = 0
	)

	row := server.mysqlDB.QueryRow(sql, req.GetAccount())
	err := row.Scan(&num)
	if err != nil {
		log.Error("mysql scan error:", err)
		return userid, protobuf.ErrCodeType_MysqlError
	}
	if num != 0 {
		return userid, protobuf.ErrCodeType_AccountExistError
	}
	updateTime := time.Now().Unix()
	sqlPrepare, err := server.mysqlDB.Prepare("INSERT INTO yfgo.player_info (account,passwd,name,autograph,icon,sex,update_time) VALUES (?,?,?,?,?,?,?)")
	if err != nil {
		log.Error("mysql scan error:", err)
		return userid, protobuf.ErrCodeType_MysqlError
	}
	result, err := sqlPrepare.Exec(req.Account, req.Passwd, req.Name, req.Autograph, req.Icon, req.Sex, updateTime)
	//result, err := server.mysqlDB.Exec(sqlPrepare, req.Account, req.Passwd, req.Name, req.Autograph, req.Icon, req.Sex, updateTime)
	if err != nil {
		log.Error("mysql scan error:", err)
		return userid, protobuf.ErrCodeType_MysqlError
	}
	lastid, err := result.LastInsertId()
	if err != nil {
		log.Error("mysql scan error:", err)
		return userid, protobuf.ErrCodeType_MysqlError
	}
	userid = int32(lastid)
	return userid, protobuf.ErrCodeType_Success
}

func GetFriendList(userid int32) *protobuf.FriendListRespose {
	res := &protobuf.FriendListRespose{}
	res.Success = protobuf.ErrCodeType_Success.Enum()
	sql := "select f.friend_userid,p.name,p.autograph,p.grade,p.icon,p.sex,p.gnum,p.gsnum,p.gfnum from yfgo.friend f join yfgo.player_info p on f.friend_userid = p.id where f.self_userid = ?"
	rows, err := server.mysqlDB.Query(sql)
	defer func() {
		if rows != nil {
			if err := rows.Close(); err != nil {
				log.Error("mysql rows close error:", err)
			}
		}
	}()
	if err != nil {
		log.Error("mysql Query error:", err)
		res.Success = protobuf.ErrCodeType_MysqlError.Enum()
		return res
	}

	for rows.Next() {
		friendInfo := &protobuf.PlayerInfo{}
		err := rows.Scan(friendInfo.Userid, friendInfo.Name, friendInfo.Autograph, friendInfo.Grade, friendInfo.Icon, friendInfo.Sex, friendInfo.Gnum, friendInfo.Gsnum, friendInfo.Gfnum)
		if err != nil {
			log.Error("mysql scan error:", err)
			res := &protobuf.FriendListRespose{}
			res.Success = protobuf.ErrCodeType_MysqlError.Enum()
			return res
		}
		player := server.playerManager.findPlayerById(friendInfo.GetUserid())
		if player != nil {
			friendInfo.Status = proto.Int32(player.status)
		}
		res.List = append(res.List, friendInfo)
	}
	return res
}

func InviteOrAccpetFriend(inviter int32, invitee int32, flag int32) int32 {
	sql := "call InviteOrAccpetFriend(?,?,?,@flag)"
	_, err := server.mysqlDB.Exec(sql, inviter, invitee, flag)
	if err != nil {
		log.Error("mysql scan error:", err)
		return -1
	}
	sql = "select @flag"
	var success int32 = -1
	err = server.mysqlDB.QueryRow(sql).Scan(&success)
	if err != nil || success == -1 {
		log.Error("mysql scan error:", err)
		return -1
	}
	return success
}

func DeleteFriend(selfUserid int32, friendUserid int32) bool {
	sql := "call DeleteFriend(?,?)"
	_, err := server.mysqlDB.Exec(sql, selfUserid, friendUserid)
	if err != nil {
		log.Error("mysql scan error:", err)
		return false
	}
	return true
}
