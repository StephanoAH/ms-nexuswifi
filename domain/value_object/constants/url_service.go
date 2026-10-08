package constants

import "fmt"

var PathSendInvoice = "https://%s:7011/api"

var SendInvoiceElectronicDefault = fmt.Sprintf(PathSendInvoice, currentEnvironment)
