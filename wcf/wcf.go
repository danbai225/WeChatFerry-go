package wcf

import (
	logs "github.com/danbai225/go-logs"
	"go.nanomsg.org/mangos/v3"
	"go.nanomsg.org/mangos/v3/protocol"
	"go.nanomsg.org/mangos/v3/protocol/pair1"
	_ "go.nanomsg.org/mangos/v3/transport/all"
	"google.golang.org/protobuf/proto"
	"strconv"
	"strings"
)

type Client struct {
	add     string
	socket  protocol.Socket
	RecvTxt bool
}

func (c *Client) conn() error {
	socket, err := pair1.NewSocket()
	if err != nil {
		return err
	}
	err = socket.Dial(c.add)
	if err != nil {
		return err
	}
	c.socket = socket
	return err
}
func (c *Client) send(data []byte) error {
	return c.socket.Send(data)
}
func (c *Client) Recv() (*Response, error) {
	msg := &Response{}
	recv, err := c.socket.Recv()
	if err != nil {
		return msg, err
	}
	err = proto.Unmarshal(recv, msg)
	return msg, err
}
func (c *Client) Close() error {
	c.DisableRecvTxt()
	return c.socket.Close()
}
func (c *Client) IsLogin() bool {
	err := c.send(genFunReq(Functions_FUNC_IS_LOGIN).build())
	if err != nil {
		logs.Err(err)
	}
	recv, err := c.Recv()
	if err != nil {
		logs.Err(err)
	}
	if recv.GetStatus() == 1 {
		return true
	}
	return false
}
func (c *Client) GetSelfWXID() string {
	err := c.send(genFunReq(Functions_FUNC_GET_SELF_WXID).build())
	if err != nil {
		logs.Err(err)
	}
	recv, err := c.Recv()
	if err != nil {
		logs.Err(err)
	}
	return recv.GetStr()
}
func (c *Client) GetMsgTypes() map[int32]string {
	err := c.send(genFunReq(Functions_FUNC_GET_MSG_TYPES).build())
	if err != nil {
		logs.Err(err)
	}
	recv, err := c.Recv()
	if err != nil {
		logs.Err(err)
	}
	return recv.GetTypes().GetTypes()
}
func (c *Client) GetContacts() []*RpcContact {
	err := c.send(genFunReq(Functions_FUNC_GET_CONTACTS).build())
	if err != nil {
		logs.Err(err)
	}
	recv, err := c.Recv()
	if err != nil {
		logs.Err(err)
	}
	return recv.GetContacts().GetContacts()
}
func (c *Client) GetDBNames() []string {
	err := c.send(genFunReq(Functions_FUNC_GET_DB_NAMES).build())
	if err != nil {
		logs.Err(err)
	}
	recv, err := c.Recv()
	if err != nil {
		logs.Err(err)
	}
	return recv.GetDbs().Names
}
func (c *Client) GetDBTables(tab string) []*DbTable {
	req := genFunReq(Functions_FUNC_GET_DB_TABLES)
	str := &Request_Str{Str: tab}
	req.Msg = str
	err := c.send(req.build())
	if err != nil {
		logs.Err(err)
	}
	recv, err := c.Recv()
	if err != nil {
		logs.Err(err)
	}
	return recv.GetTables().GetTables()
}
func (c *Client) ExecDBQuery(db, sql string) []*DbRow {
	req := genFunReq(Functions_FUNC_EXEC_DB_QUERY)
	q := Request_Query{
		Query: &DbQuery{
			Db:  db,
			Sql: sql,
		},
	}
	req.Msg = &q
	err := c.send(req.build())
	if err != nil {
		logs.Err(err)
	}
	recv, err := c.Recv()
	if err != nil {
		logs.Err(err)
	}
	return recv.GetRows().GetRows()
}

/*AcceptFriend 接收好友请求
 * 接收好友请求
 *
 * @param v3 xml.attrib["encryptusername"] // 加密的用户名
 * @param v4 xml.attrib["ticket"]   // Ticket
 * @param scene 17 // 添加方式：17 名片，30 扫码
 */
func (c *Client) AcceptFriend(v3, v4 string, scene int32) int32 {
	req := genFunReq(Functions_FUNC_ACCEPT_FRIEND)
	q := Request_V{
		V: &Verification{
			V3:    v3,
			V4:    v4,
			Scene: scene,
		}}

	req.Msg = &q
	err := c.send(req.build())
	if err != nil {
		logs.Err(err)
	}
	recv, err := c.Recv()
	if err != nil {
		logs.Err(err)
	}
	return recv.GetStatus()
}
func (c *Client) AddChatroomMembers(roomID, wxIDs string) int32 {
	req := genFunReq(Functions_FUNC_ADD_ROOM_MEMBERS)
	q := Request_M{
		M: &MemberMgmt{Roomid: roomID, Wxids: wxIDs},
	}
	req.Msg = &q
	err := c.send(req.build())
	if err != nil {
		logs.Err(err)
	}
	recv, err := c.Recv()
	if err != nil {
		logs.Err(err)
	}
	return recv.GetStatus()
}

/*
   ReceiveTransfer 接收转账
   string wxid = 1; // 转账人
   string tfid = 2; // 转账id transferid
   string taid = 3; // Transaction id
*/

func (c *Client) ReceiveTransfer(wxid, tfid, taid string) int32 {
	req := genFunReq(Functions_FUNC_RECV_TRANSFER)
	q := Request_Tf{
		Tf: &Transfer{
			Wxid: wxid,
			Tfid: tfid,
			Taid: taid,
		},
	}
	req.Msg = &q
	err := c.send(req.build())
	if err != nil {
		logs.Err(err)
	}
	recv, err := c.Recv()
	if err != nil {
		logs.Err(err)
	}
	return recv.GetStatus()
}

// RefreshPYQ 刷新朋友圈
func (c *Client) RefreshPYQ(id uint64) int32 {
	req := genFunReq(Functions_FUNC_REFRESH_PYQ)
	q := Request_Ui64{
		Ui64: id,
	}
	req.Msg = &q
	err := c.send(req.build())
	if err != nil {
		logs.Err(err)
	}
	recv, err := c.Recv()
	if err != nil {
		logs.Err(err)
	}
	return recv.GetStatus()
}

/*
DownloadAttach
下载附件（图片、视频、文件）。这方法别直接调用，下载图片使用 `download_image`。

	Args:
	    id (int): 消息中 id
	    thumb (str): 消息中的 thumb
	    extra (str): 消息中的 extra

	Returns:
	    int: 0 为成功, 其他失败。
*/
func (c *Client) DownloadAttach(id uint64, thumb string, extra string) int32 {
	req := genFunReq(Functions_FUNC_DOWNLOAD_ATTACH)
	q := Request_Att{
		Att: &AttachMsg{
			Id:    id,
			Thumb: thumb,
			Extra: extra,
		},
	}
	req.Msg = &q
	err := c.send(req.build())
	if err != nil {
		logs.Err(err)
	}
	recv, err := c.Recv()
	if err != nil {
		logs.Err(err)
	}
	return recv.GetStatus()
}

/*
GetContactInfo
通过 wxid 查询微信号昵称等信息

	Args:
	    wxid (str): 联系人 wxid

	Returns:
	    dict: {wxid, code, name, gender}
*/
func (c *Client) GetContactInfo(wxId string) *RpcContact {
	req := genFunReq(Functions_FUNC_GET_CONTACT_INFO)
	q := Request_Str{
		Str: wxId,
	}
	req.Msg = &q
	err := c.send(req.build())
	if err != nil {
		logs.Err(err)
	}
	recv, err := c.Recv()
	if err != nil {
		logs.Err(err)
	}
	contacts := recv.GetContacts().GetContacts()
	if len(contacts) > 0 {
		return contacts[0]
	}
	return nil
}

/*
Revoke
撤回消息

	Args:
	    id (int): 待撤回消息的 id

	Returns:
	    int: 1 为成功，其他失败
*/
func (c *Client) Revoke(id uint64) int32 {
	req := genFunReq(Functions_FUNC_REVOKE_MSG)
	q := Request_Ui64{
		Ui64: id,
	}
	req.Msg = &q
	err := c.send(req.build())
	if err != nil {
		logs.Err(err)
	}
	recv, err := c.Recv()
	if err != nil {
		logs.Err(err)
	}
	return recv.GetStatus()
}

// DecryptImage 解密图片 加密路径，解密路径
func (c *Client) DecryptImage(src, dst string) int32 {
	req := genFunReq(Functions_FUNC_DECRYPT_IMAGE)
	q := Request_Dec{
		Dec: &DecPath{Src: src, Dst: dst},
	}
	req.Msg = &q
	err := c.send(req.build())
	if err != nil {
		logs.Err(err)
	}
	recv, err := c.Recv()
	if err != nil {
		logs.Err(err)
	}
	return recv.GetStatus()
}

/*
ExecOCR
获取 OCR 结果。鸡肋，需要图片能自动下载；通过下载接口下载的图片无法识别。

	Args:
	    extra (str): 待识别的图片路径，消息里的 extra

	Returns:
	    str: OCR 结果
*/
func (c *Client) ExecOCR(extra string) string {
	req := genFunReq(Functions_FUNC_EXEC_OCR)
	q := Request_Str{
		Str: extra,
	}
	req.Msg = &q
	err := c.send(req.build())
	if err != nil {
		logs.Err(err)
	}
	recv, err := c.Recv()
	if err != nil {
		logs.Err(err)
	}
	return recv.GetStr()
}
func (c *Client) AddChatRoomMembers(roomId string, wxIds []string) int32 {
	req := genFunReq(Functions_FUNC_ADD_ROOM_MEMBERS)
	q := Request_M{
		M: &MemberMgmt{Roomid: roomId,
			Wxids: strings.Join(wxIds, ",")},
	}
	req.Msg = &q
	err := c.send(req.build())
	if err != nil {
		logs.Err(err)
	}
	recv, err := c.Recv()
	if err != nil {
		logs.Err(err)
	}
	return recv.GetStatus()
}
func (c *Client) DelChatRoomMembers(roomId string, wxIds []string) int32 {
	req := genFunReq(Functions_FUNC_DEL_ROOM_MEMBERS)
	q := Request_M{
		M: &MemberMgmt{Roomid: roomId,
			Wxids: strings.Join(wxIds, ",")},
	}
	req.Msg = &q
	err := c.send(req.build())
	if err != nil {
		logs.Err(err)
	}
	recv, err := c.Recv()
	if err != nil {
		logs.Err(err)
	}
	return recv.GetStatus()
}

/*
InvChatRoomMembers
邀请群成员

	Args:
	    roomid (str): 群的 id
	    wxids (str): 要邀请成员的 wxid, 多个用逗号`,`分隔

	Returns:
	    int: 1 为成功，其他失败
*/
func (c *Client) InvChatRoomMembers(roomId string, wxIds []string) int32 {
	req := genFunReq(Functions_FUNC_INV_ROOM_MEMBERS)
	q := Request_M{
		M: &MemberMgmt{Roomid: roomId,
			Wxids: strings.Join(wxIds, ",")},
	}
	req.Msg = &q
	err := c.send(req.build())
	if err != nil {
		logs.Err(err)
	}
	recv, err := c.Recv()
	if err != nil {
		logs.Err(err)
	}
	return recv.GetStatus()
}
func (c *Client) GetUserInfo() *UserInfo {
	err := c.send(genFunReq(Functions_FUNC_GET_USER_INFO).build())
	if err != nil {
		logs.Err(err)
	}
	recv, err := c.Recv()
	if err != nil {
		logs.Err(err)
	}
	return recv.GetUi()
}

// GetAudio 获取语音消息并转成 MP3
/*
Args:
            id (int): 消息中 id
            dir (str): 存放图片的目录

        Returns:
            str: 成功返回存储路径；空字符串为失败，原因见日志。
*/
func (c *Client) GetAudio(id uint64, dir string) string {
	req := genFunReq(Functions_FUNC_SEND_TXT)
	req.Msg = &Request_Am{
		Am: &AudioMsg{
			Id:  id,
			Dir: dir,
		},
	}
	err := c.send(req.build())
	if err != nil {
		logs.Err(err)
	}
	recv, err := c.Recv()
	if err != nil {
		logs.Err(err)
	}
	return recv.GetStr()
}

/*
SendTxt
@param msg:      消息内容（如果是 @ 消息则需要有跟 @ 的人数量相同的 @）
@param receiver: 消息接收人，私聊为 wxid（wxid_xxxxxxxxxxxxxx），群聊为roomid（xxxxxxxxxx@chatroom）
@param ates:    群聊时要 @ 的人（私聊时为空字符串），多个用逗号分隔。@所有人 用notify@all（必须是群主或者管理员才有权限）
*/
func (c *Client) SendTxt(msg string, receiver string, ates []string) int32 {
	req := genFunReq(Functions_FUNC_SEND_TXT)
	req.Msg = &Request_Txt{
		Txt: &TextMsg{
			Msg:      msg,
			Receiver: receiver,
			Aters:    strings.Join(ates, ","),
		},
	}
	err := c.send(req.build())
	if err != nil {
		logs.Err(err)
	}
	recv, err := c.Recv()
	if err != nil {
		logs.Err(err)
	}
	return recv.GetStatus()
}

/*
SendIMG
path 绝对路径InBot
*/
func (c *Client) SendIMG(path string, receiver string) int32 {
	req := genFunReq(Functions_FUNC_SEND_IMG)
	req.Msg = &Request_File{
		File: &PathMsg{
			Path:     path,
			Receiver: receiver,
		},
	}
	err := c.send(req.build())
	if err != nil {
		logs.Err(err)
	}
	recv, err := c.Recv()
	if err != nil {
		logs.Err(err)
	}
	return recv.GetStatus()
}

/*
SendFile
path 绝对路径InBot
*/
func (c *Client) SendFile(path string, receiver string) int32 {
	req := genFunReq(Functions_FUNC_SEND_FILE)
	req.Msg = &Request_File{
		File: &PathMsg{
			Path:     path,
			Receiver: receiver,
		},
	}
	err := c.send(req.build())
	if err != nil {
		logs.Err(err)
	}
	recv, err := c.Recv()
	if err != nil {
		logs.Err(err)
	}
	return recv.GetStatus()
}
func (c *Client) SendXml(path, content, receiver string, Type int32) int32 {
	req := genFunReq(Functions_FUNC_SEND_XML)
	req.Msg = &Request_Xml{
		Xml: &XmlMsg{
			Receiver: receiver,
			Content:  content,
			Path:     path,
			Type:     Type,
		},
	}
	err := c.send(req.build())
	if err != nil {
		logs.Err(err)
	}
	recv, err := c.Recv()
	if err != nil {
		logs.Err(err)
	}
	return recv.GetStatus()
}
func (c *Client) SendEmotion(path, receiver string) int32 {
	req := genFunReq(Functions_FUNC_SEND_EMOTION)
	req.Msg = &Request_File{
		File: &PathMsg{
			Path:     path,
			Receiver: receiver,
		},
	}
	err := c.send(req.build())
	if err != nil {
		logs.Err(err)
	}
	recv, err := c.Recv()
	if err != nil {
		logs.Err(err)
	}
	return recv.GetStatus()
}

/*
SendRichText
发送富文本消息

	卡片样式：
	    |-------------------------------------|
	    |title, 最长两行
	    |(长标题, 标题短的话这行没有)
	    |digest, 最多三行，会占位    |--------|
	    |digest, 最多三行，会占位    |thumburl|
	    |digest, 最多三行，会占位    |--------|
	    |(account logo) name
	    |-------------------------------------|
	Args:
	    name (str): 左下显示的名字
	    account (str): 填公众号 id 可以显示对应的头像（gh_ 开头的）
	    title (str): 标题，最多两行
	    digest (str): 摘要，三行
	    url (str): 点击后跳转的链接
	    thumburl (str): 缩略图的链接
	    receiver (str): 接收人, wxid 或者 roomid

	Returns:
	    int: 0 为成功，其他失败
*/
func (c *Client) SendRichText(name, account, title, digest, url, thumbUrl, receiver string) int32 {
	req := genFunReq(Functions_FUNC_SEND_RICH_TXT)
	req.Msg = &Request_Rt{
		Rt: &RichText{
			Name:     name,
			Account:  account,
			Title:    title,
			Digest:   digest,
			Url:      url,
			Thumburl: thumbUrl,
			Receiver: receiver,
		},
	}
	err := c.send(req.build())
	if err != nil {
		logs.Err(err)
	}
	recv, err := c.Recv()
	if err != nil {
		logs.Err(err)
	}
	return recv.GetStatus()
}

/*
SendPat
拍一拍群友

	Args:
	    roomid (str): 群 id
	    wxid (str): 要拍的群友的 wxid

	Returns:
	    int: 1 为成功，其他失败
*/
func (c *Client) SendPat(roomId, wxId string) int32 {
	req := genFunReq(Functions_FUNC_SEND_RICH_TXT)
	req.Msg = &Request_Pm{
		Pm: &PatMsg{
			Roomid: roomId,
			Wxid:   wxId,
		},
	}
	err := c.send(req.build())
	if err != nil {
		logs.Err(err)
	}
	recv, err := c.Recv()
	if err != nil {
		logs.Err(err)
	}
	return recv.GetStatus()
}
func (c *Client) EnableRecvTxt() int32 {
	req := genFunReq(Functions_FUNC_ENABLE_RECV_TXT)
	req.Msg = &Request_Flag{
		Flag: true,
	}
	err := c.send(req.build())
	if err != nil {
		logs.Err(err)
	}
	recv, err := c.Recv()
	if err != nil {
		logs.Err(err)
	}
	return recv.GetStatus()
}
func (c *Client) DisableRecvTxt() int32 {
	err := c.send(genFunReq(Functions_FUNC_DISABLE_RECV_TXT).build())
	if err != nil {
		logs.Err(err)
	}
	recv, err := c.Recv()
	if err != nil {
		logs.Err(err)
	}
	c.RecvTxt = false
	return recv.GetStatus()
}
func (c *Client) OnMSG(f func(msg *WxMsg)) error {
	c.RecvTxt = true
	socket, err := pair1.NewSocket()
	if err != nil {
		return err
	}
	_ = socket.SetOption(mangos.OptionRecvDeadline, 2000)
	_ = socket.SetOption(mangos.OptionSendDeadline, 2000)
	err = socket.Dial(addPort(c.add))
	if err != nil {
		return err
	}
	defer socket.Close()
	for c.RecvTxt {
		msg := &Response{}
		recv, err := socket.Recv()
		if err != nil {
			return err
		}
		_ = proto.Unmarshal(recv, msg)
		go f(msg.GetWxmsg())
	}
	return err
}
func NewWCF(add string) (*Client, error) {
	if add == "" {
		add = "tcp://127.0.0.1:10086"
	}
	client := &Client{add: add}
	err := client.conn()
	return client, err
}

type cmdMSG struct {
	*Request
}

func (c *cmdMSG) build() []byte {
	marshal, _ := proto.Marshal(c)
	return marshal
}
func genFunReq(fun Functions) *cmdMSG {
	return &cmdMSG{
		&Request{Func: fun,
			Msg: nil},
	}
}
func addPort(add string) string {
	parts := strings.Split(add, ":")
	port, _ := strconv.Atoi(parts[2])
	newPort := port + 1
	return parts[0] + ":" + parts[1] + ":" + strconv.Itoa(newPort)
}
