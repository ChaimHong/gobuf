package cssource

import (
	"log"
	"testing"
)

func TestGen(t *testing.T) {
	jsonData := `{
  "File": "struct.go",
  "Package": "main",
  "Enums": [
    
  ],
  "Structs": [
    {
      "Name": "RankPlayer",
      "Fields": [
        {
          "Name": "BestGroup",
          "Type": {
            "Kind": "Int8",
            "Name": "LoginGroup"
          }
        },
        {
          "Name": "User",
          "Type": {
            "Kind": "String"
          }
        },
        {
          "Name": "Num",
          "Type": {
            "Kind": "Int32"
          }
        },
        {
          "Name": "Distance",
          "Type": {
            "Kind": "Int64"
          }
        },
        {
          "Name": "BestAliveTime",
          "Type": {
            "Kind": "Int32"
          }
        },
        {
          "Name": "BestGolds",
          "Type": {
            "Kind": "Int32"
          }
        }
      ]
    },
    {
      "Name": "RankIn",
      "Fields": []
    },
    {
      "Name": "RankOut",
      "Fields": [
        {
          "Name": "List",
          "Type": {
            "Kind": "Array",
            "Elem": {
              "Kind": "Struct",
              "Name": "RankPlayer"
            }
          }
        },
        {
          "Name": "Rank",
          "Type": {
            "Kind": "Int32"
          }
        }
      ]
    },
    {
      "Name": "NotifyCloseSessionOut",
      "Fields": []
    },
    {
      "Name": "LoginIn",
      "Fields": [
        {
          "Name": "User",
          "Type": {
            "Kind": "Bytes"
          }
        },
        {
          "Name": "Group",
          "Type": {
            "Kind": "Int8",
            "Name": "LoginGroup"
          }
        }
      ]
    },
    {
      "Name": "LoginOut",
      "Fields": [
        {
          "Name": "Status",
          "Type": {
            "Kind": "Int8",
            "Name": "LoginStatus"
          }
        },
        {
          "Name": "PlayerId",
          "Type": {
            "Kind": "Int64"
          }
        },
        {
          "Name": "LastDistance",
          "Type": {
            "Kind": "Int64"
          }
        },
        {
          "Name": "MaxDistance",
          "Type": {
            "Kind": "Int64"
          }
        }
      ]
    },
    {
      "Name": "DistanceIn",
      "Fields": [
        {
          "Name": "Distance",
          "Type": {
            "Kind": "Int64"
          }
        },
        {
          "Name": "AliveTime",
          "Type": {
            "Kind": "Int32"
          }
        },
        {
          "Name": "Golds",
          "Type": {
            "Kind": "Int32"
          }
        }
      ]
    },
    {
      "Name": "DistanceOut",
      "Fields": []
    }
  ]
}`
	gencode, err := Gen([]byte(jsonData))
	log.Printf("%s %v", gencode, err)
}
