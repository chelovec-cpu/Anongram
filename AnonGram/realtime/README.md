# WebRTC

WebRTC медиа-канал, WebSocket/WSS — signaling. Для production нужен STUN/TURN (например coturn), потому что P2P работает не через каждый NAT.
Сигналинг: offer -> server -> answer, плюс ICE candidates.
