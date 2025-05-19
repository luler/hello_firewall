package logic

import (
	"fmt"
	"gin_base/app/helper/db_helper"
	"gin_base/app/helper/exception_helper"
	"gin_base/app/model"
	"os/exec"
	"strings"
	"sync"
)

// IPTablesManager 管理iptables规则
type IPTablesManager struct {
	Chain     string
	Signature string // 规则标识
}

// NewIPTablesManager 创建新的iptables管理器
func NewIPTablesManager(chain string, signature string) *IPTablesManager {
	return &IPTablesManager{
		Chain:     chain,
		Signature: signature,
	}
}

// InitChain 初始化自定义链
func (m *IPTablesManager) InitChain() error {
	// 使用 iptables-save 检查链是否存在，比 iptables -L 更高效
	checkCmd := exec.Command("bash", "-c", "iptables-save | grep -w '^:"+m.Chain+"'")
	if err := checkCmd.Run(); err != nil {
		// 创建新链
		createCmd := exec.Command("iptables", "-N", m.Chain)
		if err := createCmd.Run(); err != nil {
			return fmt.Errorf("创建链失败: %v", err)
		}
	}

	// 无论链是否已经存在，都确保它在INPUT链的第一位
	// 首先删除现有的规则（如果存在）
	delCmd := exec.Command("bash", "-c", "iptables -D INPUT -j "+m.Chain+" -m comment --comment \""+m.Signature+"\" 2>/dev/null || true")
	delCmd.Run() // 忽略错误，因为规则可能不存在
	// 然后将链添加到INPUT链的第一位，并添加标识
	linkCmd := exec.Command("iptables", "-I", "INPUT", "1", "-j", m.Chain,
		"-m", "comment", "--comment", m.Signature)
	if err := linkCmd.Run(); err != nil {
		return fmt.Errorf("将规则链链接到INPUT链首位失败: %v", err)
	}

	//DOCKER-USER链兼容
	// 添加到DOCKER-USER链（如果存在）
	checkDockerUserCmd := exec.Command("iptables", "-L", "DOCKER-USER")
	if checkDockerUserCmd.Run() == nil { // 如果DOCKER-USER链存在

		delDockerUserCmd := exec.Command("bash", "-c", "iptables -D DOCKER-USER -j "+m.Chain+" -m comment --comment \""+m.Signature+"\" 2>/dev/null || true")
		delDockerUserCmd.Run() // 忽略错误，因为规则可能不存在

		// 然后将链添加到INPUT链的第一位，并添加标识
		linkDockerUserCmd := exec.Command("iptables", "-I", "DOCKER-USER", "1", "-j", m.Chain,
			"-m", "comment", "--comment", m.Signature)
		if err := linkDockerUserCmd.Run(); err != nil {
			return fmt.Errorf("将规则链链接到DockerUser链首位失败: %v", err)
		}
	}

	return nil
}

// ClearAllRules 清除所有带有特定标识的规则
func (m *IPTablesManager) ClearAllRules() error {
	// 获取所有规则
	cmd := exec.Command("iptables", "-S", m.Chain)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	lines := strings.Split(string(output), "\n")

	// 删除所有包含标识的规则
	for _, line := range lines {
		if strings.Contains(line, m.Signature) {
			// 将 -A 替换为 -D 来构建删除命令
			deleteRule := strings.Replace(line, "-A", "-D", 1)
			parts := strings.Fields(deleteRule)

			if len(parts) > 0 {
				deleteCmd := exec.Command("iptables", parts...)
				deleteCmd.Run() // 忽略错误，尝试删除尽可能多的规则
			}
		}
	}

	return nil
}

// ApplyRule 应用单条规则
func (m *IPTablesManager) ApplyRule(rule *model.IPRule) error {
	if rule.Status == 0 {
		return nil // 跳过禁用的规则
	}

	args := []string{"-A", m.Chain, "-s", rule.IP}

	// 添加协议和端口（如果有）
	if rule.Protocol != "" {
		args = append(args, "-p", rule.Protocol)

		if rule.Port > 0 {
			args = append(args, "--dport", fmt.Sprintf("%d", rule.Port))
		}
	}

	// 添加动作和标识
	args = append(args, "-j", "DROP", "-m", "comment", "--comment", m.Signature)

	// 执行命令
	cmd := exec.Command("iptables", args...)
	return cmd.Run()
}

// RebuildRules 重建所有规则
func (m *IPTablesManager) RebuildRules() error {
	//删除后查询所有规则，重置iptables规则
	var rules []*model.IPRule
	db_helper.Db().Where("status = 1").Find(&rules)
	// 先清除所有规则
	if err := m.ClearAllRules(); err != nil {
		return err
	}

	// 应用所有启用的规则
	for _, rule := range rules {
		if err := m.ApplyRule(rule); err != nil {
			return err
		}
	}

	return nil
}

var iptables *IPTablesManager
var go_once sync.Once

// 获取单例
func GetIPTablesManager() *IPTablesManager {
	go_once.Do(func() {
		iptables = NewIPTablesManager("HELLO-FIREWALL", "managed-by-hello-firewall")
		// 初始化自定义链
		if err := iptables.InitChain(); err != nil {
			exception_helper.CommonException(fmt.Sprintf("初始化iptables链失败: %v", err))
		}
	})
	return iptables
}
