package capability

import (
	"fmt"
	"os"
	"time"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func (ca CapAudit) Report(opts ...ReportOption) error {

	for _, opt := range opts {

		err := opt.report(ca)
		if err != nil {
			logger.Debugf("Unable to generate report for %T", opt, "Error: %s", err)
		}

	}
	return nil
}

type ReportOption interface {
	report(ca CapAudit) error
}

// Simple print option implmentation for operator install
type OperatorInstallRptOptionPrint struct{}

func (OperatorInstallRptOptionPrint) report(ca CapAudit) error {

	fmt.Println()
	fmt.Println("Operator Install Report:")
	fmt.Println("-----------------------------------------")
	fmt.Printf("Report Date: %s\n", time.Now())
	fmt.Printf("OpenShift Version: %s\n", ca.OcpVersion)
	fmt.Printf("Package Name: %s\n", ca.Subscription.Package)
	fmt.Printf("Channel: %s\n", ca.Subscription.Channel)
	fmt.Printf("Catalog Source: %s\n", ca.Subscription.CatalogSource)
	fmt.Printf("Install Mode: %s\n", ca.Subscription.InstallModeType)

	if !ca.CsvTimeout {
		fmt.Printf("Result: %s\n", ca.Csv.Status.Phase)
	} else {
		fmt.Println("Result: timeout")
	}

	fmt.Printf("Message: %s\n", ca.Csv.Status.Message)
	fmt.Printf("Reason: %s\n", ca.Csv.Status.Reason)
	fmt.Println("-----------------------------------------")

	return nil
}

// Simple file option implementation for operator install
type OperatorInstallRptOptionFile struct {
	FilePath string
}

func (opt OperatorInstallRptOptionFile) report(ca CapAudit) error {

	file, err := os.OpenFile(opt.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		file.Close()
		return err
	}
	defer file.Close()

	if !ca.CsvTimeout {

		file.WriteString("{\"level\":\"info\",\"message\":\"" + string(ca.Csv.Status.Phase) + "\",\"package\":\"" + ca.Subscription.Package + "\",\"channel\":\"" + ca.Subscription.Channel + "\",\"installmode\":\"" + string(ca.Subscription.InstallModeType) + "\"}\n")
	} else {

		file.WriteString("{\"level\":\"info\",\"message\":\"" + "timeout" + "\",\"package\":\"" + ca.Subscription.Package + "\",\"channel\":\"" + ca.Subscription.Channel + "\",\"installmode\":\"" + string(ca.Subscription.InstallModeType) + "\"}\n")
	}

	return nil
}

// Simple print option implmentation for operand install
type OperandInstallRptOptionPrint struct{}

func (OperandInstallRptOptionPrint) report(ca CapAudit) error {

	operand := &unstructured.Unstructured{Object: ca.CustomResources[0]}

	fmt.Println()
	fmt.Println("Operand Install Report:")
	fmt.Println("-----------------------------------------")
	fmt.Printf("Report Date: %s\n", time.Now())
	fmt.Printf("OpenShift Version: %s\n", ca.OcpVersion)
	fmt.Printf("Package Name: %s\n", ca.Subscription.Package)
	fmt.Printf("Operand Kind: %s\n", operand.GetKind())
	fmt.Printf("Operand Name: %s\n", operand.GetName())

	if len(ca.Operands) > 0 {
		fmt.Println("Operand Creation: Succeeded")
	} else {
		fmt.Println("Operand Creation: Failed")
	}
	fmt.Println("-----------------------------------------")

	return nil
}

// Simple file option implementation for operand install
type OperandInstallRptOptionFile struct {
	FilePath string
}

func (opt OperandInstallRptOptionFile) report(ca CapAudit) error {

	operand := &unstructured.Unstructured{Object: ca.CustomResources[0]}

	file, err := os.OpenFile(opt.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		file.Close()
		return err
	}
	defer file.Close()

	if len(ca.Operands) > 0 {

		file.WriteString("{\"package\":\"" + ca.Subscription.Package + "\", \"Operand Kind\": \"" + operand.GetKind() + "\", \"Operand Name\": \"" + operand.GetName() + "\",\"message\":\"" + "created" + "\"}\n")
	} else {

		file.WriteString("{\"package\":\"" + ca.Subscription.Package + "\", \"Operand Kind\": \"" + operand.GetKind() + "\", \"Operand Name\": \"" + operand.GetName() + "\",\"message\":\"" + "failed" + "\"}\n")
	}

	return nil
}
