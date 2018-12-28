// Copyright 2018 The QOS Authors

package commands

import (
	"errors"
	"fmt"

	"github.com/QOSGroup/qmoon/db"
	"github.com/QOSGroup/qmoon/service"
	"github.com/QOSGroup/qmoon/types"
	"github.com/QOSGroup/qmoon/utils"
	"github.com/spf13/cobra"
)

// NodeCmd 数据库初始化命令
var NodeCmd = &cobra.Command{
	Use:   "node",
	Short: "node 管理",
}

var createNodeCmd = &cobra.Command{
	Use:   "create",
	Short: "添加node节点",
	RunE:  createNode,
}

var queryNodeCmd = &cobra.Command{
	Use:   "query",
	Short: "查询node节点",
	RunE:  queryNode,
}

var updateNodeCmd = &cobra.Command{
	Use:   "update nodeName",
	Short: "更新node节点",
	RunE:  updateNode,
}

var deleteNodeCmd = &cobra.Command{
	Use:   "delete nodeName",
	Short: "删除node节点",
	RunE:  deleteNode,
}

var (
	nodeName string
	nodeUrl  string
	nodeType string
)

func init() {
	createNodeCmd.PersistentFlags().StringVar(&nodeName, "nodeName", "", "the name of node")
	createNodeCmd.PersistentFlags().StringVar(&nodeUrl, "nodeUrl", "", "the url of node")
	createNodeCmd.PersistentFlags().StringVar(&nodeType, "nodeType", "", fmt.Sprintf("节点类型:%s, %s",
		types.NodeTypeQOS, types.NodeTypeQSC))

	updateNodeCmd.PersistentFlags().StringVar(&nodeName, "nodeName", "", "the name of node")
	updateNodeCmd.PersistentFlags().StringVar(&nodeUrl, "nodeUrl", "", "the url of node")

	queryNodeCmd.PersistentFlags().StringVar(&nodeName, "nodeName", "", "the name of node")

	registerFlagsDb(createNodeCmd)
	registerFlagsDb(queryNodeCmd)
	registerFlagsDb(updateNodeCmd)
	registerFlagsDb(deleteNodeCmd)

	NodeCmd.AddCommand(createNodeCmd, queryNodeCmd, deleteNodeCmd)
}

func createNode(cmd *cobra.Command, args []string) error {
	err := db.InitDb(config.DB, logger)
	if err != nil {
		return err
	}

	if nodeName == "" {
		return errors.New("nodeName 不能为空")
	}

	if nodeUrl == "" {
		return errors.New("nodeUrl 不能为空")
	}

	if nodeType == "" {
		return errors.New("nodeType 不能为空")
	}

	if !types.CheckNodeType(nodeType) {
		return errors.New("nodeType 不支持")
	}

	err = service.CreateNode(nodeName, nodeUrl, nodeType, "", nil)
	if err != nil {
		return err
	}

	return nil
}

func queryNode(cmd *cobra.Command, args []string) error {
	err := db.InitDb(config.DB, logger)
	if err != nil {
		return err
	}

	headers := []string{"name", "chain_id", "url"}
	var datas [][]string
	if nodeName != "" {
		res, err := service.GetNodeByName(nodeName)
		if err != nil {
			return err
		}

		datas = append(datas, []string{res.Name, res.ChanID, res.BaseURL})
		utils.PrintTable(cmd.OutOrStdout(), headers, datas)
	} else {
		res, err := service.AllNodes()
		if err != nil {
			return err
		}

		for _, v := range res {
			datas = append(datas, []string{v.Name, v.ChanID, v.BaseURL})
		}

		utils.PrintTable(cmd.OutOrStdout(), headers, datas)
	}

	return nil
}

func updateNode(cmd *cobra.Command, args []string) error {

	return nil
}

func deleteNode(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("需要参数nodeName")
	}
	name := args[0]

	err := db.InitDb(config.DB, logger)
	if err != nil {
		return err
	}

	err = service.DeleteNodeByName(name)

	return err
}
