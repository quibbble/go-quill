# Events

## Core events that are hookable
- damage_unit x
- heal_unit x
- move_unit x
- attack_unit x
- kill_unit x
- place_unit x
- summon_unit x
- add_item_to_unit x
- remove_item_from_unit x
- play_card x
- discard_card x
- draw_card x
- burn_card x
- sack_card x
- drain_mana x
- gain_mana x
- drain_base_mana x
- gain_base_mana x
- recycle_deck x
- end_turn x

## Non hookable events
- damage_units x
- heal_units x
- add_trait_to_card x
- remove_trait_from_card x
- swap_units x
- end_game x
- refresh_movement x
- cooldown x

## Traits to implement
- poison X - COMPLETE
- berserk - COMPLETE
- recode X - COMPLETE
- buff STAT X - COMPLETE
- debuff STAT X - COMPLETE
- execute - COMPLETE
- shield X - COMPLETE
- ward X - COMPLETE
- thief - COMPLETE
- purity - COMPLETE
- pillage EVENT - COMPLETE
- battle cry EVENT - COMPLETE
- death cry EVENT - COMPLETE
- gift TRAIT - COMPLETE
- lobber - COMPLETE
- friends CHOOSE TRAIT - 
- enemies CHOOSE TRAIT - 
- spiky X - COMPLETE
- enrage EVENT - COMPLETE
- assassin X - COMPLETE

friends + enemies
Need keep track of who had before and who had after
Need to check both unit(s) in event and all other units
- On place
- On spawn
- On move
- On swap
- On kill
- On add trait
- On remove trait
