package auth

import (
	"errors"
	"fmt"
	"github.com/scofieldpeng/config-go"
	"strings"
)

type (
	RequestMethod string // 请求方式
	ApiInfo       struct {
		Method RequestMethod //请求方式
		Auth   bool          // 授权方式
	} // Api接口信息

	ApiPath struct {
		Path      string                    // api路径
		ApiMethod map[RequestMethod]ApiInfo // 请求方式的接口信息
	} // 某API路径下的api信息

	ApiName struct {
		Path string
		ApiInfo
	} //  某个api名称下的授权信息
)

const (
	GetMethod    RequestMethod = "GET"    // GET请求
	PostMethod   RequestMethod = "POST"   // POST请求
	PutMethod    RequestMethod = "PUT"    // PUT请求
	DeleteMethod RequestMethod = "DELETE" // DELETE请求
)

var (
	apiPathList map[string]map[string]ApiPath // api路径map列表,第一维key为大版本号,第二维的key为api接口的路径
	apiNameList map[string]map[string]ApiName // api名称map列表,第一维key为大版本号,第二维的key为api接口的路径
)

// Init 初始化
func Init() error {
	apiPathList = make(map[string]map[string]ApiPath)
	apiNameList = make(map[string]map[string]ApiName)

	apiConfigFile := config.Config("api")
	for sectionName, section := range apiConfigFile {
		splitName := strings.Split(sectionName, "_")
		if len(splitName) != 2 {

			errors.New(fmt.Sprintf("api.ini文件中配置[%s]命名不规范,请设置为`大版本号_接口名`", sectionName))
		}
		if _, ok := apiPathList[splitName[0]]; !ok {
			apiPathList[splitName[0]] = make(map[string]ApiPath)
		}
		if _, ok := apiNameList[splitName[0]]; !ok {
			apiNameList[splitName[0]] = make(map[string]ApiName)
		}

		// 获取每个api接口配置信息
		path := section["path"]
		if path == "" {
			return errors.New(fmt.Sprintf("api.ini文件中配置[%s]中缺少path配置", sectionName))
		}
		methodStr := section["method"]
		if methodStr == "" {
			return errors.New(fmt.Sprintf("api.ini文件中配置[%s]中缺少method配置", sectionName))
		}
		method := RequestMethod(strings.ToUpper(methodStr))
		if method != GetMethod && method != PostMethod && method != PutMethod && method != DeleteMethod {
			return errors.New(fmt.Sprintf("api.ini文件中配置[%s]中method配置不正确", sectionName))
		}
		auth := config.Bool(section["auth"], true)

		// 写入apiNameList
		if _, ok := apiNameList[splitName[0]][splitName[1]]; !ok {
			apiNameList[splitName[0]][splitName[1]] = ApiName{
				Path: path,
				ApiInfo: ApiInfo{
					Method: method,
					Auth:   auth,
				},
			}
		}
		// 写入apiPathList
		if _, ok := apiPathList[splitName[0]][path]; !ok {
			apiPathList[splitName[0]][path] = ApiPath{
				Path:      path,
				ApiMethod: make(map[RequestMethod]ApiInfo),
			}
		}
		apiPathList[splitName[0]][path].ApiMethod[method] = ApiInfo{
			Method: method,
			Auth:   auth,
		}
	}

	return nil
}

// ApiNeedToAuth api是否需要授权,接口名称,请求方式即可,返回bool值
func ApiNeedAuth(path string, method RequestMethod) bool {
	version := getVersion(path)
	if version == "" {
		return true
	}
	if pathInfos, ok := apiPathList[version]; !ok {
		//fmt.Println("没找到version")
		return true
	} else {
		if pathInfo, ok := pathInfos[path]; !ok {
			//fmt.Println("没找到path:",path)
			return true
		} else {
			if methodInfo, ok := pathInfo.ApiMethod[method]; !ok {
				//fmt.Println("没找到path对应的method,path:",path,",method:",method)
				return true
			} else {
				return methodInfo.Auth
			}
		}
	}
}

// getVersion 获取请求的版本号
func getVersion(path string) string {
	// add slash in the path end
	if string([]byte(path)[len(path)-1:]) != "" {
		path = path + "/"
	}
	return strings.Split(strings.Split(path, "/api")[1], "/")[1]
}
