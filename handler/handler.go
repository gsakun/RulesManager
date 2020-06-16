package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var AlertGroupInfomap map[int]string = make(map[int]string)

type Ruler struct {
	Groups []Group `yaml:"groups"`
}

type Group struct {
	Name  string `yaml:"name"`
	Rules []Rule `yaml:"rules"`
}

type Rule struct {
	Alert       string            `yaml:"alert"`
	Annotations map[string]string `yaml:"annotations"`
	Expr        string            `yaml:"expr"`
	For         string            `yaml:"for"`
	Labels      map[string]string `yaml:"labels"`
}

func Queryalertgroup(db *sql.DB) {
	sql := "select id, alert_group_name from alert_group"
	rows, err := db.Query(sql)
	if err != nil {
		log.Errorf("Query alert_group table Failed")
	} else {
		defer rows.Close()
		for rows.Next() {
			var (
				id               int
				alert_group_name string
			)

			err = rows.Scan(&id, &alert_group_name)
			if err != nil {
				log.Errorf("ERROR: %v", err)
				continue
			}
			AlertGroupInfomap[id] = alert_group_name
		}
	}
	//log.Infof("sync alert_group table success %v", AlertGroupInfomap)
}

func HandlerRule(db *sql.DB, requiredtype, path string) error {
	sql := fmt.Sprintf("select alert_rule_name,expr,alert_trigger_time,labels,annotations,alert_group_id,state from alert_rule where state!=0 and alert_group_id in (%s)", requiredtype)
	rows, err := db.Query(sql)
	if err != nil {
		log.Errorf("Query data_center table Failed")
		return err
	} else {
		defer rows.Close()
		var errnum int
		for rows.Next() {
			var (
				alert_rule_name    string
				expr               string
				alert_trigger_time string
				labels             string
				annotations        string
				alert_group_id     int
				state              int
			)

			err = rows.Scan(&alert_rule_name, &expr, &alert_trigger_time, &labels, &annotations, &alert_group_id, &state)
			if err != nil {
				log.Errorf("ERROR: %v", err)
				errnum++
				continue
			}
			if state == 1 || state == 2 {
				err := generateYaml(alert_rule_name, expr, alert_trigger_time, labels, annotations, path, alert_group_id)
				if err != nil {
					fmt.Errorf("Generate yaml file failed for %s-%s,errinfo %v", AlertGroupInfomap[alert_group_id], alert_rule_name, err)
					errnum++
					continue
				}
				stmt, err := db.Prepare(`UPDATE alert_rule set state=? where alert_rule_name=? and alert_group_id=? `)
				if err != nil {
					log.Errorf("Update prepare err %s-%s %v", AlertGroupInfomap[alert_group_id], alert_rule_name, err)
					errnum++
					continue
				}
				_, err = stmt.Exec(0, alert_rule_name, alert_group_id)
				if err != nil {
					log.Errorf("Update exec err %s-%s %v", AlertGroupInfomap[alert_group_id], alert_rule_name, err)
					errnum++
					continue
				}
			}
			if state == 3 {
				filename := fmt.Sprintf("%s%s___%s.yml", path, AlertGroupInfomap[alert_group_id], alert_rule_name)
				err := os.Remove(filename)
				if err != nil {
					log.Errorf("Remove file %s failed, errinfo %v", filename, err)
					errnum++
					continue
				}
				stmt, err := db.Prepare("DELETE FROM alert_rule WHERE alert_group_id=? and alert_rule_name=?")
				if err != nil {
					errnum++
					log.Errorf("DELETE prepare failed %v", err)
					continue
				}
				_, err = stmt.Exec(alert_group_id, alert_rule_name)
				if err != nil {
					errnum++
					log.Errorf("DELETE exec failed %v", err)
					continue
				}
			}
		}
		if errnum != 0 {
			return fmt.Errorf("Encountered several errors")
		}
	}
	return nil
}

func generateYaml(alert_rule_name, expr, alert_trigger_time, labels, annotations, path string, alert_group_id int) error {
	var rule Rule
	rule.Alert = alert_rule_name
	if labels != "" {
		var rulelabels map[string]string = make(map[string]string)
		var m map[string]interface{}
		err := json.Unmarshal([]byte(labels), &m)
		if err != nil {
			return fmt.Errorf("Parse Label failed errinfo %v", err)
		}
		if len(m) != 0 {
			for key, value := range m {
				switch value.(type) {
				case string:
					rulelabels[key] = value.(string)
				case int:
					rulelabels[key] = strconv.Itoa(value.(int))
				case float64:
					rulelabels[key] = strconv.FormatFloat(value.(float64), 'E', -1, 64)
				default:
					continue
				}
			}
		}
		if len(rulelabels) != 0 {
			rule.Labels = rulelabels
		}
	}
	if annotations != "" {
		var ruleannotations map[string]string = make(map[string]string)
		var m map[string]interface{}
		err := json.Unmarshal([]byte(annotations), &m)
		if err != nil {
			return fmt.Errorf("Parse Annotations failed errinfo %v", err)
		}
		if len(m) != 0 {
			for key, value := range m {
				switch value.(type) {
				case string:
					ruleannotations[key] = value.(string)
				case int:
					ruleannotations[key] = strconv.Itoa(value.(int))
				case float64:
					ruleannotations[key] = strconv.FormatFloat(value.(float64), 'E', -1, 64)
				default:
					continue
				}
			}
		}
		if len(ruleannotations) != 0 {
			rule.Annotations = ruleannotations
		}
	}
	rule.Expr = expr
	rule.For = alert_trigger_time
	var group Group
	group.Name = AlertGroupInfomap[alert_group_id]
	group.Rules = append(group.Rules, rule)
	var ruler Ruler
	ruler.Groups = append(ruler.Groups, group)
	d, err := yaml.Marshal(&ruler)
	if err != nil {
		return fmt.Errorf("Marshal yaml failed errinfo %v", err)
	}
	filename := fmt.Sprintf("%s%s___%s.yml", path, AlertGroupInfomap[alert_group_id], alert_rule_name)
	err = writeToFile(filename, string(d))
	if err != nil {
		return fmt.Errorf("WriteToFile Failed errinfo %v", err)
	}
	return nil
}

func writeToFile(fileName string, content string) error {
	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		log.Errorf("file create failed. err: %v", err)
	} else {
		// offset
		//os.Truncate(filename, 0) //clear
		n, _ := f.Seek(0, os.SEEK_END)
		_, err = f.WriteAt([]byte(content), n)
		defer f.Close()
	}
	return err
}
