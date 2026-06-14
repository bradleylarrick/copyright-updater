#!/bin/bash
REPO=~bradley/.m2/repository
PGM_CLASSES=target/classes

cp=${PGM_CLASSES}:$REPO/com/google/guava/guava/33.6.0-jre/guava-33.6.0-jre.jar

java -cp "$cp" org.natuna.datagen.Launcher "$@"
