#!/bin/bash

result=$(./sysclient accountAdd 0x98dd8c8F3d14D68C7b3965B2f002480a2E540C1B 0x206075f758f210d571293ca7a2be51f09930a6c7947eea5638fee5b3e2922635)

if grep -q "succesfully" <<< "$result";
then
	echo $result
	echo "koble"
else
	echo "poulous"
fi
