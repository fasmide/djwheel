#include "FastLED.h"
#include <Encoder.h>

#define NUM_LEDS 26
#define DATA_PIN 5
CRGB leds[NUM_LEDS];

Encoder encoder(2, 3);
long oldPosition  = -999;


// This function sets up the ledsand tells the controller about them
void setup() {
    Serial.println("Basic Encoder Test:");
   	delay(2000);
    FastLED.addLeds<WS2812B, DATA_PIN, GRB>(leds, NUM_LEDS);
}

void loop() {

   // output buffer every 16 miliseconds
   EVERY_N_MILLISECONDS(16) {
       FastLED.show();
       fadeToBlackBy(leds, NUM_LEDS, 10);
   }

   long newPosition = encoder.read();
   if (newPosition != oldPosition) {
     oldPosition = newPosition;
     Serial.println(newPosition);
     if (oldPosition >= 0 && oldPosition <= 500) {
       leds[map(oldPosition, 0, 500, 0, 25)] = CRGB::White;
     }
  }
}
