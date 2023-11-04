package helpers

import (
	"context"
	"dash/database"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func FileNameWithoutExtension(fileName string) string {
	ext := filepath.Ext(fileName)

	return fileName[:len(fileName)-len(ext)]
}

func CreateSegment(inputName string, outputName string, workDir string, id string) {
	s := "ffmpeg -i " + inputName + " -c:v libx264 -profile:v baseline -level 3.0 -b:v 3M -c:a aac -strict -2 -f dash " + outputName

	fmt.Println(outputName)

	args := strings.Split(s, " ")
	cmd := exec.Command(args[0], args[1:]...)

	cmd.Dir = workDir

	lol, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Println("The error:", err)
		fmt.Println("The bytes:", string(lol))
		return
	}

	collection := database.Client.Database("videos").Collection("dash")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"isdash": true}}
	_, updateErr := collection.UpdateOne(ctx, filter, update)
	if updateErr != nil {
		fmt.Println("Error while updating isdash in mongodb:", updateErr.Error())
		return
	}

}
