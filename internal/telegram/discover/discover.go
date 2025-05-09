package discover

import (
	"fmt"
	"log"
	"os/exec"
)

func ExecuteDiscoverNodeExporter() (string, error) {
	scriptPath := "/home/discover/discover_prometheus_nodes.sh"

	cmd := exec.Command("bash", scriptPath)
	output, err := cmd.CombinedOutput()

	if err != nil {
		log.Printf("Error al ejecutar el script %s: %v, salida: %s", scriptPath, err, output)
		return "", fmt.Errorf("error al ejecutar el script %s: %w, salida: %s", scriptPath, err, output)
	}

	return string(output), nil
}

func ExecuteDiscoverPortExporter() (string, error) {
	scriptPath := "/home/discover/discover_prometheus_ports.sh"

	cmd := exec.Command("bash", scriptPath)
	output, err := cmd.CombinedOutput()

	if err != nil {
		log.Printf("Error al ejecutar el script %s: %v, salida: %s", scriptPath, err, output)
		return "", fmt.Errorf("error al ejecutar el script %s: %w, salida: %s", scriptPath, err, output)
	}

	return string(output), nil
}
