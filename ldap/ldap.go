package ldap

import (
	"fmt"
	"regexp"

	"github.com/breadysimon/goless/logging"
	"github.com/go-ldap/ldap"
)

var log *logging.Logger = logging.GetLogger()

// CheckAuth 用来验证LDAP用户
func CheckAuth(host string, port int, bindUserName, bindUserPassword string, baseDn, filterTemplate string, username, password string) (dn string, err error) {
	//  host: ldap server name or ip
	// 	port: generally is 389
	// 	bindUserName: 用来获取查询权限的 bind 用户。通常是DN格式
	// 	bindUserPassword:
	// 	baseDn: 从这个节点开始搜索
	// 	filterTemplate: 查询模板, 必须括号包围,如"(&(objectClass=user)(sAMAccountName=%s))"
	// 	username, password: 要验证的用户和密码

	url := fmt.Sprintf("%s:%d", host, port)
	var l *ldap.Conn
	if l, err = ldap.Dial("tcp", url); err == nil {
		log.Debug("connected:", host)

		defer l.Close()

		// 先用我们的 bind 账号给 bind 上去
		if err = l.Bind(bindUserName, bindUserPassword); err == nil {
			log.Debug("search user bond:", bindUserName)
			// 构造查询请求
			searchRequest := ldap.NewSearchRequest(
				// 这里是 basedn，我们将从这个节点开始搜索
				baseDn,
				// 这里几个参数分别是 scope, derefAliases, sizeLimit, timeLimit,  typesOnly
				// 详情可以参考 RFC4511 中的定义
				ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
				// 这里是 LDAP 查询的 Filter
				fmt.Sprintf(filterTemplate, username),
				// 这里是查询返回的属性，以数组形式提供。如果为空则会返回所有的属性
				[]string{"dn"},
				nil,
			)
			log.Debug(searchRequest)
			// 搜索返回的是一个数组
			var sr *ldap.SearchResult
			if sr, err = l.Search(searchRequest); err == nil {
				log.Debug("search result:", sr)
				// 如果没有数据返回或者超过1条数据返回，这对于用户认证而言都是不允许的。
				// 前这意味着没有查到用户，后者意味着存在重复数据
				if len(sr.Entries) != 1 {
					log.Info("User does not exist or too many entries returned: ", username)
					return "", nil
				}

				// 如果没有意外，那么我们就可以获取用户的实际 DN 了
				userdn := sr.Entries[0].DN
				log.Debug("User DN:", userdn)

				// Bind as the user to verify their password
				// 拿这个 dn 和他的密码去做 bind 验证
				if err = l.Bind(userdn, password); err == nil {
					return userdn, nil
				}

				// Rebind as the read only user for any further queries
				// 如果后续还需要做其他操作，那么使用最初的 bind 账号重新 bind 回来。恢复初始权限。
				// err = l.Bind(bindusername, bindpassword)
				// if err != nil {
				// 	logger.Fatal(err)
				// }
			}
		}
	}
	log.Error(err)
	return "", err
}

func GetCN(dn string) string {
	r := regexp.MustCompile("CN=(.*?),")
	matches := r.FindStringSubmatch(dn)
	log.Debug("matches:", matches)
	return matches[1]
}
