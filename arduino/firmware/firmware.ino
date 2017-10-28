// FastLED from http://fastled.io/
#include "FastLED.h"
// Encoder from https://www.pjrc.com/teensy/td_libs_Encoder.html
#include <Encoder.h>

// initialize vars for leds
#define NUM_LEDS 26
#define DATA_PIN 5
CRGB leds[NUM_LEDS];

// initialize encoder (pins 2 and 3)
Encoder encoder(2, 3);
long oldPosition  = -999;

void setup() {
  FastLED.addLeds<WS2812B, DATA_PIN, GRB>(leds, NUM_LEDS);
}

void loop() {
  // If there is data in our buffer, read a whole strip
  if (Serial.available()) {
    Serial.readBytes( (char*)leds, NUM_LEDS * 3);
  }

  // output led buffer every 16 miliseconds
  EVERY_N_MILLISECONDS(16) {
    FastLED.show();
    fadeToBlackBy(leds, NUM_LEDS, 10);
  }

  // read encoder position, if it moved send data.
  // TODO: figureout somekind of debouncing (maybe it must move at least 2?)
  long newPosition = encoder.read();
  if (newPosition != oldPosition) {
    oldPosition = newPosition;
    Serial.println(newPosition);
  }
}
