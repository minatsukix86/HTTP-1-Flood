Description de l'outil d'attaque HTTP Flood Ce référentiel contient une implémentation basée sur Go d'un outil d'attaque HTTP Flood conçu pour simuler des attaques par déni de service distribué (DDoS) sur une cible spécifiée. L'outil est principalement destiné à des fins éducatives pour aider les développeurs et les professionnels de la sécurité réseau à comprendre les mécanismes des attaques DDoS, leurs implications et comment s'en défendre.

Avertissement : cet outil est destiné à un usage éducatif uniquement. L'utilisation non autorisée contre tout système sans autorisation explicite est illégale et contraire à l'éthique. Obtenez toujours l'autorisation appropriée avant d'effectuer toute forme de test.

Fonctionnalités

Prise en charge HTTPS : effectue des attaques par inondation via HTTPS à l'aide de serveurs proxy.
Multithreading : lance plusieurs threads pour augmenter l'efficacité de l'attaque.
Prise en charge du proxy : lit les configurations de proxy à partir d'un fichier pour un anonymat amélioré.
Gestion des réponses : lit et imprime la réponse du serveur cible à des fins de surveillance.

Comment ça marche

Établit une connexion : l'outil se connecte à un serveur proxy spécifié pour initier la requête vers la cible.
Envoie une requête HTTP CONNECT : il envoie une requête CONNECT pour canaliser la requête HTTP via le proxy.
Inonde les requêtes : l'outil envoie plusieurs requêtes HTTP GET au débit spécifié en utilisant plusieurs threads pour une génération de charge efficace.
Lit les réponses : il capture et imprime les réponses du serveur pour surveiller les effets de l'attaque.
