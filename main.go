package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	mqtt "git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"
	"github.com/spf13/cobra"
	"github.com/teemow/mqweather/wunderground"
	"github.com/teemow/mqweather/wunderground/features"
)

var (
	globalFlags struct {
		debug   bool
		verbose bool
	}

	mainFlags struct {
		Host     string
		Port     int
		ApiKey   string
		Station  string
		Interval int
	}

	mainCmd = &cobra.Command{
		Use:   "mqweather",
		Short: "Publish weather via MQTT",
		Long:  "Send MQTT messages to a broker. Using the weather underground API on a Raspberry Pi",
		Run:   mainRun,
	}

	projectVersion string
	projectBuild   string
)

func init() {
	mainCmd.PersistentFlags().BoolVarP(&globalFlags.debug, "debug", "d", false, "Print debug output")
	mainCmd.PersistentFlags().BoolVarP(&globalFlags.verbose, "verbose", "v", false, "Print verbose output")
	mainCmd.PersistentFlags().StringVar(&mainFlags.Host, "host", "localhost", "MQTT host")
	mainCmd.PersistentFlags().IntVar(&mainFlags.Port, "port", 1883, "MQTT port")
	mainCmd.PersistentFlags().StringVar(&mainFlags.ApiKey, "api-key", "", "Weather Underground API key")
	mainCmd.PersistentFlags().StringVar(&mainFlags.Station, "station", "", "Weather Underground Station ID")
	mainCmd.PersistentFlags().IntVar(&mainFlags.Interval, "interval", 60, "Interval in seconds")
}

func assert(err error) {
	if err != nil {
		if globalFlags.debug {
			fmt.Printf("%#v\n", err)
			os.Exit(1)
		} else {
			log.Fatal(err)
		}
	}
}

func getConditions(w *wunderground.Wunderground) (*features.Conditions, error) {
	conditions, err := w.Conditions(mainFlags.Station)
	if err != nil {
		return &features.Conditions{}, err
	}

	return conditions.Condition, nil
}

func watchWeather(client *mqtt.Client, w *wunderground.Wunderground) {
	for {
		conditions, err := getConditions(w)
		if err == nil {
			// temperature
			topic := fmt.Sprintf("mqweather/%s/temperature", mainFlags.Station)
			payload := fmt.Sprintf("%d", int(conditions.TempC*1000))

			if token := client.Publish(topic, 0, false, payload); token.Wait() && token.Error() != nil {
				fmt.Printf("Failed to send message: %v\n", token.Error())
			}

			// dewpoint
			topic = fmt.Sprintf("mqweather/%s/dewpoint", mainFlags.Station)
			payload = fmt.Sprintf("%d", int(conditions.DewpointC*1000))

			if token := client.Publish(topic, 0, false, payload); token.Wait() && token.Error() != nil {
				fmt.Printf("Failed to send message: %v\n", token.Error())
			}

			// wind
			topic = fmt.Sprintf("mqweather/%s/wind", mainFlags.Station)
			payload = fmt.Sprintf("%d", int(conditions.WindKPH*1000))

			if token := client.Publish(topic, 0, false, payload); token.Wait() && token.Error() != nil {
				fmt.Printf("Failed to send message: %v\n", token.Error())
			}

			// pressure
			topic = fmt.Sprintf("mqweather/%s/pressure", mainFlags.Station)
			payload = fmt.Sprintf("%s", conditions.PressureMB)

			if token := client.Publish(topic, 0, false, payload); token.Wait() && token.Error() != nil {
				fmt.Printf("Failed to send message: %v\n", token.Error())
			}

			// humidity
			topic = fmt.Sprintf("mqweather/%s/humidity", mainFlags.Station)
			payload = fmt.Sprintf("%s", strings.Replace(conditions.RelativeHumidity, "%", "", 1))

			if token := client.Publish(topic, 0, false, payload); token.Wait() && token.Error() != nil {
				fmt.Printf("Failed to send message: %v\n", token.Error())
			}

		} else {
			fmt.Println("Could not read weather", err)
		}
		time.Sleep(time.Duration(mainFlags.Interval) * time.Second)
	}
}

func mainRun(cmd *cobra.Command, args []string) {
	// mqtt
	opts := mqtt.NewClientOptions().AddBroker(fmt.Sprintf("tcp://%s:%d", mainFlags.Host, mainFlags.Port))
	opts.SetClientID(fmt.Sprintf("mqweather-%s", mainFlags.Station))

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Println(token.Error())
	}

	conf := wunderground.DefaultConfig(mainFlags.ApiKey)
	w := wunderground.New(conf)

	go watchWeather(client, w)

	// Handle SIGINT and SIGTERM.
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)

	client.Disconnect(250)
}

func main() {
	mainCmd.AddCommand(versionCmd)

	mainCmd.Execute()
}
