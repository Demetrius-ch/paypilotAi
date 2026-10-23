# Cahier des Charges — PayPilot AI

**Version:** 1.0
**Date:** 23 Octobre 2026
**Projet:** Hackathon PayPal + AI (Devpost)

---

## 1. Présentation du Projet

### 1.1 Contexte

PayPilot AI est une application web développée dans le cadre du hackathon "Build What's Next with PayPal and AI" organisé par PayPal sur Devpost. Le projet vise à démontrer une intégration significative de PayPal et de l'intelligence artificielle.

### 1.2 Vision

> « Comprends tes paiements. Automatise ton activité. Garde le contrôle. »

PayPilot AI est un assistant intelligent destiné aux freelances et petites entreprises pour gérer leurs commandes, suivre leurs paiements PayPal et prendre de meilleures décisions grâce à un agent IA conversationnel.

### 1.3 Public Cible

| Segment | Description |
|---------|-------------|
| **Freelances** | Développeurs, designers, consultants qui facturent des prestations |
| **Micro-entrepreneurs** | Vendeurs de services ou produits digitaux |
| **Petites agences** | Équipes de 1-5 personnes gérant plusieurs clients |

### 1.4 Problèmes Résolus

| Problème | Impact | Solution PayPilot |
|----------|--------|-------------------|
| Suivi manuel des paiements | Perte de temps, erreurs | Vue centralisée + statuts automatiques |
| Relances clients oubliées | Trésorerie fragilisée | Agent IA génère des rappels |
| Pas de vue d'ensemble | Décisions difficiles | Dashboard + analyse IA |
| Création de liens de paiement | Processus répétitif | Génération en langage naturel |
| Données dispersées | Inefficacité | Centralisation dans une app |

---

## 2. Exigences du Hackathon

### 2.1 Exigences Obligatoires

- [ ] Intégration significative d'au moins une technologie PayPal (API, SDK)
- [ ] Intégration significative de l'IA dans l'expérience ou la fonctionnalité
- [ ] Prototype fonctionnel (pas un mockup statique)
- [ ] Documentation suffisante pour les juges
- [ ] Dépôt de code public (GitHub) avec licence open source
- [ ] Vidéo de démonstration < 3 minutes
- [ ] URL de démo hébergée OU instructions claires pour exécuter le projet

### 2.2 Critères d'Évaluation

1. **Créativité** — Originalité de l'idée
2. **Impact** — Valeur réelle pour les utilisateurs
3. **Exécution** — Qualité technique et finition
4. **Utilisation de PayPal** — Intégration pertinente et profonde
5. **Utilisation de l'IA** — IA utile, pas superficielle

---

## 3. Spécifications Fonctionnelles

### 3.1 Vue d'Ensemble des Fonctionnalités

```
┌─────────────────────────────────────────────────────────────────────────┐
│                           PAYPILOT AI                                   │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────────────┐  │
│  │   FONC. A       │  │   FONC. B       │  │   FONC. C               │  │
│  │   Paiements     │  │   Agent IA      │  │   Analyse               │  │
│  │   PayPal        │  │   Conversationnel│  │   Commerciale          │  │
│  └────────┬────────┘  └────────┬────────┘  └────────────┬────────────┘  │
│           │                    │                        │               │
│           └────────────────────┼────────────────────────┘               │
│                                │                                        │
│                    ┌───────────┴───────────┐                            │
│                    │      FONC. D          │                            │
│                    │  Validation Humaine   │                            │
│                    │  (Human-in-the-loop)  │                            │
│                    └───────────────────────┘                            │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

### 3.2 Fonctionnalité A — Paiements PayPal

#### Description
Gestion complète du cycle de vie des commandes avec intégration PayPal Sandbox.

#### User Stories

| ID | En tant que... | Je veux... | Afin de... |
|----|----------------|------------|------------|
| A1 | Utilisateur | Créer une commande avec montant et description | Facturer un client |
| A2 | Utilisateur | Générer un lien de paiement PayPal | Envoyer au client pour paiement |
| A3 | Utilisateur | Voir le statut de paiement en temps réel | Savoir si j'ai été payé |
| A4 | Utilisateur | Lister toutes mes commandes | Avoir une vue d'ensemble |
| A5 | Utilisateur | Filtrer par statut (payé, en attente, annulé) | Trouver rapidement une commande |

#### Flux de Paiement

```
┌──────────────┐     ┌──────────────┐     ┌──────────────┐     ┌──────────────┐
│   Création   │────▶│   Lien       │────▶│   Paiement   │────▶│   Confirmation│
│   Commande   │     │   PayPal     │     │   Client     │     │   Webhook    │
└──────────────┘     └──────────────┘     └──────────────┘     └──────────────┘
       │                    │                    │                    │
       ▼                    ▼                    ▼                    ▼
  POST /orders      PayPal API          Sandbox PayPal        POST /webhook
  Status: pending   order_id créé       Popup paiement        Status: paid
```

#### Règles Métier

- Le montant minimum d'une commande est 1.00 USD
- Les devises supportées: USD, EUR (extensible)
- Une commande ne peut pas être modifiée après création du lien PayPal
- Le statut est mis à jour automatiquement via webhook PayPal

### 3.3 Fonctionnalité B — Agent IA Conversationnel

#### Description
Interface de chat permettant à l'utilisateur d'interagir en langage naturel avec ses données et d'effectuer des actions.

#### User Stories

| ID | En tant que... | Je veux... | Afin de... |
|----|----------------|------------|------------|
| B1 | Utilisateur | Poser des questions sur mes commandes | Obtenir des infos rapidement |
| B2 | Utilisateur | Demander un résumé de mes revenus | Comprendre ma situation |
| B3 | Utilisateur | Demander de créer un lien de paiement | Gagner du temps |
| B4 | Utilisateur | Demander de préparer un message de relance | Communiquer efficacement |
| B5 | Utilisateur | Confirmer ou refuser les actions proposées | Garder le contrôle |

#### Architecture de l'Agent

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              AGENT IA                                       │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ┌─────────────────┐                                                        │
│  │  Message        │                                                        │
│  │  Utilisateur    │                                                        │
│  └────────┬────────┘                                                        │
│           │                                                                 │
│           ▼                                                                 │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │                         CLAUDE API                                  │    │
│  │                    (avec Function Calling)                          │    │
│  │                                                                     │    │
│  │   System Prompt:                                                    │    │
│  │   "Tu es PayPilot, un assistant pour gérer les paiements..."       │    │
│  │                                                                     │    │
│  │   Tools disponibles:                                                │    │
│  │   - list_orders                                                     │    │
│  │   - get_order_details                                               │    │
│  │   - create_payment_link (nécessite confirmation)                    │    │
│  │   - get_revenue_summary                                             │    │
│  │   - draft_reminder_message                                          │    │
│  │                                                                     │    │
│  └────────┬────────────────────────────────────────────────────────────┘    │
│           │                                                                 │
│           ▼                                                                 │
│  ┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐        │
│  │  Tool Call      │────▶│  Exécution      │────▶│  Résultat       │        │
│  │  (choix IA)     │     │  Backend        │     │  → Claude       │        │
│  └─────────────────┘     └─────────────────┘     └─────────────────┘        │
│                                                         │                   │
│                                                         ▼                   │
│                                                ┌─────────────────┐          │
│                                                │  Réponse        │          │
│                                                │  Naturelle      │          │
│                                                └─────────────────┘          │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

#### Tools de l'Agent

| Tool | Description | Paramètres | Confirmation |
|------|-------------|------------|--------------|
| `list_orders` | Liste les commandes | `status?`, `period?`, `limit?` | Non |
| `get_order_details` | Détails d'une commande | `order_id` | Non |
| `create_payment_link` | Crée un lien PayPal | `order_id` | **Oui** |
| `get_revenue_summary` | Résumé des revenus | `period` | Non |
| `draft_reminder_message` | Prépare une relance | `order_id`, `tone?` | Non |
| `send_reminder` | Envoie la relance | `order_id`, `message` | **Oui** |

#### Exemples d'Interactions

```
Utilisateur: "Montre-moi les commandes en attente"
Agent: [Tool: list_orders(status="pending")]
       "Tu as 3 commandes en attente:
        1. Site web pour Jean Dupont - 500€ (créée il y a 5 jours)
        2. Logo pour Marie Martin - 150€ (créée il y a 2 jours)
        3. Consultation pour Pierre Durand - 200€ (créée aujourd'hui)"

Utilisateur: "Génère un lien de paiement pour la commande de Jean"
Agent: [Tool: create_payment_link(order_id="...")]
       "Je vais créer un lien de paiement de 500€ pour Jean Dupont.

        ⚠️ Action en attente de confirmation:
        [Confirmer] [Annuler]"

Utilisateur: [Confirme]
Agent: "Lien créé! Voici le lien de paiement:
        https://paypal.com/pay/..."
```

### 3.4 Fonctionnalité C — Analyse Commerciale

#### Description
Dashboard et insights sur l'activité commerciale basés sur les données de commandes et paiements.

#### User Stories

| ID | En tant que... | Je veux... | Afin de... |
|----|----------------|------------|------------|
| C1 | Utilisateur | Voir mes revenus du mois | Suivre ma performance |
| C2 | Utilisateur | Voir le nombre de commandes par statut | Identifier les blocages |
| C3 | Utilisateur | Comparer avec le mois précédent | Voir l'évolution |
| C4 | Utilisateur | Identifier les commandes en retard | Agir rapidement |

#### Dashboard

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                            DASHBOARD                                        │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐              │
│  │  Revenus Mois   │  │  Commandes      │  │  En Attente     │              │
│  │                 │  │  Payées         │  │                 │              │
│  │    2,450 €      │  │      8          │  │    3 (850€)     │              │
│  │    +15% ▲       │  │                 │  │                 │              │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘              │
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │                     Revenus (6 derniers mois)                       │    │
│  │                                                                     │    │
│  │   3000 │                                           ╭───╮            │    │
│  │        │                               ╭───╮       │   │            │    │
│  │   2000 │               ╭───╮   ╭───╮   │   │       │   │            │    │
│  │        │       ╭───╮   │   │   │   │   │   │       │   │            │    │
│  │   1000 │       │   │   │   │   │   │   │   │       │   │            │    │
│  │        │───────┴───┴───┴───┴───┴───┴───┴───┴───────┴───┴────        │    │
│  │          Mai   Juin  Juil  Août  Sept  Oct                          │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │                     Commandes en retard                             │    │
│  │                                                                     │    │
│  │   ⚠️ Jean Dupont - Site web - 500€ - 5 jours d'attente              │    │
│  │      [Relancer] [Voir détails]                                      │    │
│  │                                                                     │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 3.5 Fonctionnalité D — Validation Humaine (Human-in-the-loop)

#### Description
Système de confirmation pour les actions sensibles avant leur exécution.

#### Actions Nécessitant Confirmation

| Action | Raison | Interface |
|--------|--------|-----------|
| Créer un lien de paiement | Impact financier | Modal de confirmation |
| Envoyer un email/message | Communication externe | Prévisualisation + confirmation |
| Annuler une commande | Irréversible | Double confirmation |

#### Flux de Validation

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         FLUX DE VALIDATION                                  │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  1. Agent propose une action                                                │
│     ┌─────────────────────────────────────────────────────────────────┐     │
│     │  "Je vais créer un lien de paiement de 500€ pour Jean Dupont"  │     │
│     └─────────────────────────────────────────────────────────────────┘     │
│                                    │                                        │
│                                    ▼                                        │
│  2. Interface affiche les détails + boutons                                 │
│     ┌─────────────────────────────────────────────────────────────────┐     │
│     │  ┌─────────────────────────────────────────────────────────┐   │     │
│     │  │  ACTION EN ATTENTE                                      │   │     │
│     │  ├─────────────────────────────────────────────────────────┤   │     │
│     │  │  Type: Création lien de paiement                        │   │     │
│     │  │  Commande: Site web pour Jean Dupont                    │   │     │
│     │  │  Montant: 500€                                          │   │     │
│     │  │  Client: jean.dupont@email.com                          │   │     │
│     │  ├─────────────────────────────────────────────────────────┤   │     │
│     │  │         [✓ Confirmer]        [✗ Annuler]                │   │     │
│     │  └─────────────────────────────────────────────────────────┘   │     │
│     └─────────────────────────────────────────────────────────────────┘     │
│                                    │                                        │
│                         ┌──────────┴──────────┐                             │
│                         │                     │                             │
│                         ▼                     ▼                             │
│  3a. Confirmation                   3b. Annulation                          │
│      → Exécution                        → Message de refus                  │
│      → Notification succès              → Retour au chat                    │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 4. Architecture Technique

### 4.1 Stack Technologique

| Couche | Technologie | Justification |
|--------|-------------|---------------|
| **Frontend** | React 18 + TypeScript | Robuste, typé, écosystème riche |
| **Styling** | Tailwind CSS | Rapide, responsive, personnalisable |
| **State** | Zustand | Léger, simple, performant |
| **Build** | Vite | Rapide, moderne |
| **Backend** | Go 1.21+ | Performant, typage fort, bon pour les API |
| **Router** | Chi ou Gin | Léger, standard |
| **Database** | SQLite | Simple, pas de serveur, suffisant pour MVP |
| **IA** | Claude API (Anthropic) | Function calling, français, performant |
| **Paiements** | PayPal REST API v2 | Exigence hackathon |

### 4.2 Architecture Globale

```
┌─────────────────────────────────────────────────────────────────────────────────────┐
│                                   ARCHITECTURE                                       │
├─────────────────────────────────────────────────────────────────────────────────────┤
│                                                                                     │
│                              ┌─────────────────────┐                                │
│                              │      NAVIGATEUR     │                                │
│                              │    (React SPA)      │                                │
│                              └──────────┬──────────┘                                │
│                                         │                                           │
│                                         │ HTTPS                                     │
│                                         │                                           │
│  ┌──────────────────────────────────────┼──────────────────────────────────────┐    │
│  │                                      │                                      │    │
│  │                            ┌─────────▼─────────┐                            │    │
│  │                            │   REVERSE PROXY   │                            │    │
│  │                            │   (Nginx/Caddy)   │                            │    │
│  │                            └─────────┬─────────┘                            │    │
│  │                                      │                                      │    │
│  │              ┌───────────────────────┼───────────────────────┐              │    │
│  │              │                       │                       │              │    │
│  │              ▼                       ▼                       ▼              │    │
│  │    ┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐      │    │
│  │    │   /api/*        │     │   /static/*     │     │   /webhook/*   │      │    │
│  │    │   Backend Go    │     │   Frontend      │     │   PayPal       │      │    │
│  │    └────────┬────────┘     │   (Build)       │     │   Callbacks    │      │    │
│  │             │              └─────────────────┘     └────────┬───────┘      │    │
│  │             │                                               │              │    │
│  │             │              ┌─────────────────────────────────┘              │    │
│  │             │              │                                               │    │
│  │             ▼              ▼                                               │    │
│  │    ┌───────────────────────────────────────────────────────────────────┐   │    │
│  │    │                        BACKEND GO                                 │   │    │
│  │    │                                                                   │   │    │
│  │    │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐               │   │    │
│  │    │  │  Handlers   │  │  Services   │  │  Database   │               │   │    │
│  │    │  │  (HTTP)     │──│  (Business) │──│  (SQLite)   │               │   │    │
│  │    │  └─────────────┘  └──────┬──────┘  └─────────────┘               │   │    │
│  │    │                          │                                       │   │    │
│  │    │         ┌────────────────┼────────────────┐                      │   │    │
│  │    │         │                │                │                      │   │    │
│  │    │         ▼                ▼                ▼                      │   │    │
│  │    │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐               │   │    │
│  │    │  │  PayPal     │  │  Claude     │  │  Email      │               │   │    │
│  │    │  │  Client     │  │  Client     │  │  Client     │               │   │    │
│  │    │  └──────┬──────┘  └──────┬──────┘  └─────────────┘               │   │    │
│  │    │         │                │                                       │   │    │
│  │    └─────────┼────────────────┼───────────────────────────────────────┘   │    │
│  │              │                │                                           │    │
│  └──────────────┼────────────────┼───────────────────────────────────────────┘    │
│                 │                │                                                │
│                 ▼                ▼                                                │
│        ┌─────────────────┐  ┌─────────────────┐                                   │
│        │   PayPal API    │  │   Claude API    │                                   │
│        │   (Sandbox)     │  │   (Anthropic)   │                                   │
│        └─────────────────┘  └─────────────────┘                                   │
│                                                                                   │
└───────────────────────────────────────────────────────────────────────────────────┘
```

### 4.3 Architecture Backend (Go)

```
┌─────────────────────────────────────────────────────────────────────────────────────┐
│                              BACKEND GO - ARCHITECTURE                              │
├─────────────────────────────────────────────────────────────────────────────────────┤
│                                                                                     │
│  ┌─────────────────────────────────────────────────────────────────────────────┐    │
│  │                              HTTP LAYER                                     │    │
│  │                                                                             │    │
│  │   ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐       │    │
│  │   │  Middleware │  │  Middleware │  │  Middleware │  │  Middleware │       │    │
│  │   │  CORS       │─▶│  Logger     │─▶│  Recovery   │─▶│  Auth (opt) │       │    │
│  │   └─────────────┘  └─────────────┘  └─────────────┘  └──────┬──────┘       │    │
│  │                                                              │              │    │
│  │                                                              ▼              │    │
│  │   ┌─────────────────────────────────────────────────────────────────────┐  │    │
│  │   │                           ROUTER                                    │  │    │
│  │   │                                                                     │  │    │
│  │   │   /api/orders/*     ──────────────▶  OrderHandler                   │  │    │
│  │   │   /api/chat         ──────────────▶  ChatHandler                    │  │    │
│  │   │   /api/analytics/*  ──────────────▶  AnalyticsHandler               │  │    │
│  │   │   /api/paypal/*     ──────────────▶  PayPalHandler                  │  │    │
│  │   │   /api/actions/*    ──────────────▶  ActionHandler                  │  │    │
│  │   │                                                                     │  │    │
│  │   └─────────────────────────────────────────────────────────────────────┘  │    │
│  │                                                                             │    │
│  └─────────────────────────────────────────────────────────────────────────────┘    │
│                                         │                                           │
│                                         ▼                                           │
│  ┌─────────────────────────────────────────────────────────────────────────────┐    │
│  │                            SERVICE LAYER                                    │    │
│  │                                                                             │    │
│  │   ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐             │    │
│  │   │  OrderService   │  │  AIService      │  │  PayPalService  │             │    │
│  │   │                 │  │                 │  │                 │             │    │
│  │   │  - Create       │  │  - ProcessMsg   │  │  - CreateOrder  │             │    │
│  │   │  - GetByID      │  │  - ExecuteTool  │  │  - CapturePaymt │             │    │
│  │   │  - List         │  │  - ListTools    │  │  - GetStatus    │             │    │
│  │   │  - UpdateStatus │  │                 │  │  - HandleWebhook│             │    │
│  │   └────────┬────────┘  └────────┬────────┘  └────────┬────────┘             │    │
│  │            │                    │                    │                      │    │
│  │            │                    │                    │                      │    │
│  │   ┌────────▼────────┐  ┌────────▼────────┐  ┌────────▼────────┐             │    │
│  │   │  AnalyticsServ  │  │  ActionService  │  │  NotifyService  │             │    │
│  │   │                 │  │                 │  │                 │             │    │
│  │   │  - GetSummary   │  │  - Create       │  │  - Email        │             │    │
│  │   │  - GetTrends    │  │  - Confirm      │  │  - (future)     │             │    │
│  │   │  - GetOverdue   │  │  - Reject       │  │                 │             │    │
│  │   └─────────────────┘  └─────────────────┘  └─────────────────┘             │    │
│  │                                                                             │    │
│  └─────────────────────────────────────────────────────────────────────────────┘    │
│                                         │                                           │
│                                         ▼                                           │
│  ┌─────────────────────────────────────────────────────────────────────────────┐    │
│  │                          REPOSITORY LAYER                                   │    │
│  │                                                                             │    │
│  │   ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐             │    │
│  │   │  OrderRepo      │  │  ChatRepo       │  │  ActionRepo     │             │    │
│  │   │  (orders table) │  │  (messages)     │  │  (pending_acts) │             │    │
│  │   └────────┬────────┘  └────────┬────────┘  └────────┬────────┘             │    │
│  │            │                    │                    │                      │    │
│  │            └────────────────────┼────────────────────┘                      │    │
│  │                                 │                                           │    │
│  │                                 ▼                                           │    │
│  │                    ┌─────────────────────────┐                              │    │
│  │                    │        SQLite DB        │                              │    │
│  │                    │      (paypilot.db)      │                              │    │
│  │                    └─────────────────────────┘                              │    │
│  │                                                                             │    │
│  └─────────────────────────────────────────────────────────────────────────────┘    │
│                                                                                     │
└─────────────────────────────────────────────────────────────────────────────────────┘
```

### 4.4 Structure des Packages Go

```
backend/
├── cmd/
│   └── server/
│       └── main.go                 # Point d'entrée, initialisation
│
├── internal/
│   ├── config/
│   │   └── config.go               # Chargement configuration (env)
│   │
│   ├── database/
│   │   ├── database.go             # Connexion SQLite
│   │   └── migrations.go           # Création des tables
│   │
│   ├── models/
│   │   ├── order.go                # struct Order
│   │   ├── message.go              # struct ChatMessage
│   │   ├── action.go               # struct PendingAction
│   │   └── analytics.go            # struct Summary, Trend
│   │
│   ├── repository/
│   │   ├── order_repo.go           # CRUD orders
│   │   ├── message_repo.go         # CRUD messages
│   │   └── action_repo.go          # CRUD pending_actions
│   │
│   ├── services/
│   │   ├── order_service.go        # Logique métier commandes
│   │   ├── analytics_service.go    # Calculs et agrégations
│   │   ├── action_service.go       # Gestion actions en attente
│   │   │
│   │   ├── paypal/
│   │   │   ├── client.go           # Client HTTP PayPal
│   │   │   ├── orders.go           # Créer/capturer orders
│   │   │   └── webhook.go          # Traitement webhooks
│   │   │
│   │   └── ai/
│   │       ├── agent.go            # Orchestration agent
│   │       ├── client.go           # Client Claude API
│   │       ├── tools.go            # Définition des tools
│   │       └── executor.go         # Exécution des tools
│   │
│   ├── handlers/
│   │   ├── order_handler.go        # Endpoints /api/orders
│   │   ├── chat_handler.go         # Endpoint /api/chat
│   │   ├── analytics_handler.go    # Endpoints /api/analytics
│   │   ├── paypal_handler.go       # Endpoints /api/paypal
│   │   └── action_handler.go       # Endpoints /api/actions
│   │
│   └── middleware/
│       ├── cors.go                 # CORS headers
│       ├── logger.go               # Request logging
│       └── recovery.go             # Panic recovery
│
├── pkg/
│   └── response/
│       └── json.go                 # Helpers réponses JSON
│
├── go.mod
├── go.sum
└── Makefile
```

### 4.5 Architecture Frontend (React)

```
┌─────────────────────────────────────────────────────────────────────────────────────┐
│                           FRONTEND REACT - ARCHITECTURE                             │
├─────────────────────────────────────────────────────────────────────────────────────┤
│                                                                                     │
│  ┌─────────────────────────────────────────────────────────────────────────────┐    │
│  │                              APP SHELL                                      │    │
│  │                                                                             │    │
│  │   ┌─────────────────────────────────────────────────────────────────────┐   │    │
│  │   │                          Layout                                     │   │    │
│  │   │   ┌─────────┐  ┌─────────────────────────────────────────────────┐  │   │    │
│  │   │   │         │  │                                                 │  │   │    │
│  │   │   │  Side   │  │                   MAIN CONTENT                  │  │   │    │
│  │   │   │  bar    │  │                                                 │  │   │    │
│  │   │   │         │  │   ┌─────────────────────────────────────────┐   │  │   │    │
│  │   │   │  - Home │  │   │                                         │   │  │   │    │
│  │   │   │  - Orders│ │   │              PAGE CONTENT                │   │  │   │    │
│  │   │   │  - Chat │  │   │           (React Router)                 │   │  │   │    │
│  │   │   │         │  │   │                                         │   │  │   │    │
│  │   │   │         │  │   └─────────────────────────────────────────┘   │  │   │    │
│  │   │   │         │  │                                                 │  │   │    │
│  │   │   └─────────┘  └─────────────────────────────────────────────────┘  │   │    │
│  │   └─────────────────────────────────────────────────────────────────────┘   │    │
│  │                                                                             │    │
│  └─────────────────────────────────────────────────────────────────────────────┘    │
│                                                                                     │
│  ┌─────────────────────────────────────────────────────────────────────────────┐    │
│  │                              PAGES                                          │    │
│  │                                                                             │    │
│  │   /                        /orders                    /orders/:id           │    │
│  │   ┌─────────────┐          ┌─────────────┐           ┌─────────────┐        │    │
│  │   │  Dashboard  │          │  OrderList  │           │ OrderDetail │        │    │
│  │   │  Page       │          │  Page       │           │ Page        │        │    │
│  │   └─────────────┘          └─────────────┘           └─────────────┘        │    │
│  │                                                                             │    │
│  │   /orders/new              /chat                                            │    │
│  │   ┌─────────────┐          ┌─────────────┐                                  │    │
│  │   │  NewOrder   │          │  Chat       │                                  │    │
│  │   │  Page       │          │  Page       │                                  │    │
│  │   └─────────────┘          └─────────────┘                                  │    │
│  │                                                                             │    │
│  └─────────────────────────────────────────────────────────────────────────────┘    │
│                                                                                     │
│  ┌─────────────────────────────────────────────────────────────────────────────┐    │
│  │                            COMPONENTS                                       │    │
│  │                                                                             │    │
│  │   ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐             │    │
│  │   │  ui/            │  │  dashboard/     │  │  orders/        │             │    │
│  │   │  - Button       │  │  - StatCard     │  │  - OrderCard    │             │    │
│  │   │  - Card         │  │  - RevenueChart │  │  - OrderForm    │             │    │
│  │   │  - Input        │  │  - OverdueList  │  │  - OrderStatus  │             │    │
│  │   │  - Modal        │  │                 │  │  - PayPalButton │             │    │
│  │   │  - Spinner      │  │                 │  │                 │             │    │
│  │   └─────────────────┘  └─────────────────┘  └─────────────────┘             │    │
│  │                                                                             │    │
│  │   ┌─────────────────┐  ┌─────────────────┐                                  │    │
│  │   │  chat/          │  │  common/        │                                  │    │
│  │   │  - ChatWindow   │  │  - Sidebar      │                                  │    │
│  │   │  - MessageList  │  │  - Header       │                                  │    │
│  │   │  - MessageInput │  │  - Layout       │                                  │    │
│  │   │  - ActionCard   │  │                 │                                  │    │
│  │   │  - ConfirmModal │  │                 │                                  │    │
│  │   └─────────────────┘  └─────────────────┘                                  │    │
│  │                                                                             │    │
│  └─────────────────────────────────────────────────────────────────────────────┘    │
│                                                                                     │
│  ┌─────────────────────────────────────────────────────────────────────────────┐    │
│  │                          STATE MANAGEMENT (Zustand)                         │    │
│  │                                                                             │    │
│  │   ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐             │    │
│  │   │  useOrderStore  │  │  useChatStore   │  │  useUIStore     │             │    │
│  │   │                 │  │                 │  │                 │             │    │
│  │   │  - orders[]     │  │  - messages[]   │  │  - sidebarOpen  │             │    │
│  │   │  - loading      │  │  - pending      │  │  - theme        │             │    │
│  │   │  - fetchOrders  │  │  - sendMessage  │  │  - notifications│             │    │
│  │   │  - createOrder  │  │  - confirmAction│  │                 │             │    │
│  │   └─────────────────┘  └─────────────────┘  └─────────────────┘             │    │
│  │                                                                             │    │
│  └─────────────────────────────────────────────────────────────────────────────┘    │
│                                                                                     │
│  ┌─────────────────────────────────────────────────────────────────────────────┐    │
│  │                           API SERVICE LAYER                                 │    │
│  │                                                                             │    │
│  │   ┌─────────────────────────────────────────────────────────────────────┐   │    │
│  │   │  services/api.ts                                                    │   │    │
│  │   │                                                                     │   │    │
│  │   │  - api.orders.list()                                                │   │    │
│  │   │  - api.orders.create(data)                                          │   │    │
│  │   │  - api.orders.getById(id)                                           │   │    │
│  │   │  - api.chat.sendMessage(content)                                    │   │    │
│  │   │  - api.actions.confirm(id)                                          │   │    │
│  │   │  - api.analytics.getSummary()                                       │   │    │
│  │   │                                                                     │   │    │
│  │   └─────────────────────────────────────────────────────────────────────┘   │    │
│  │                                                                             │    │
│  └─────────────────────────────────────────────────────────────────────────────┘    │
│                                                                                     │
└─────────────────────────────────────────────────────────────────────────────────────┘
```

### 4.6 Structure des Fichiers Frontend

```
frontend/
├── public/
│   └── favicon.ico
│
├── src/
│   ├── components/
│   │   ├── ui/
│   │   │   ├── Button.tsx
│   │   │   ├── Card.tsx
│   │   │   ├── Input.tsx
│   │   │   ├── Modal.tsx
│   │   │   ├── Spinner.tsx
│   │   │   └── index.ts
│   │   │
│   │   ├── common/
│   │   │   ├── Layout.tsx
│   │   │   ├── Sidebar.tsx
│   │   │   └── Header.tsx
│   │   │
│   │   ├── dashboard/
│   │   │   ├── StatCard.tsx
│   │   │   ├── RevenueChart.tsx
│   │   │   └── OverdueList.tsx
│   │   │
│   │   ├── orders/
│   │   │   ├── OrderCard.tsx
│   │   │   ├── OrderForm.tsx
│   │   │   ├── OrderStatus.tsx
│   │   │   └── PayPalButton.tsx
│   │   │
│   │   └── chat/
│   │       ├── ChatWindow.tsx
│   │       ├── MessageList.tsx
│   │       ├── MessageInput.tsx
│   │       ├── ActionCard.tsx
│   │       └── ConfirmModal.tsx
│   │
│   ├── pages/
│   │   ├── DashboardPage.tsx
│   │   ├── OrderListPage.tsx
│   │   ├── OrderDetailPage.tsx
│   │   ├── NewOrderPage.tsx
│   │   └── ChatPage.tsx
│   │
│   ├── hooks/
│   │   ├── useOrders.ts
│   │   └── useChat.ts
│   │
│   ├── services/
│   │   └── api.ts
│   │
│   ├── store/
│   │   ├── orderStore.ts
│   │   ├── chatStore.ts
│   │   └── uiStore.ts
│   │
│   ├── types/
│   │   ├── order.ts
│   │   ├── chat.ts
│   │   └── api.ts
│   │
│   ├── utils/
│   │   ├── format.ts
│   │   └── constants.ts
│   │
│   ├── App.tsx
│   ├── main.tsx
│   └── index.css
│
├── package.json
├── tsconfig.json
├── tailwind.config.js
├── postcss.config.js
└── vite.config.ts
```

---

## 5. Modèle de Données

### 5.1 Schéma de Base de Données

```
┌─────────────────────────────────────────────────────────────────────────────────────┐
│                              SCHÉMA BASE DE DONNÉES                                 │
├─────────────────────────────────────────────────────────────────────────────────────┤
│                                                                                     │
│  ┌─────────────────────────────────────────────────────────────────────────────┐    │
│  │                              ORDERS                                         │    │
│  ├─────────────────────────────────────────────────────────────────────────────┤    │
│  │  id              TEXT PRIMARY KEY                                           │    │
│  │  client_name     TEXT NOT NULL                                              │    │
│  │  client_email    TEXT                                                       │    │
│  │  description     TEXT                                                       │    │
│  │  amount          REAL NOT NULL                                              │    │
│  │  currency        TEXT DEFAULT 'USD'                                         │    │
│  │  status          TEXT DEFAULT 'pending'  -- pending | paid | cancelled     │    │
│  │  paypal_order_id TEXT                                                       │    │
│  │  payment_link    TEXT                                                       │    │
│  │  created_at      DATETIME DEFAULT CURRENT_TIMESTAMP                         │    │
│  │  updated_at      DATETIME DEFAULT CURRENT_TIMESTAMP                         │    │
│  │  paid_at         DATETIME                                                   │    │
│  └─────────────────────────────────────────────────────────────────────────────┘    │
│                                         │                                           │
│                                         │ 1:N                                       │
│                                         ▼                                           │
│  ┌─────────────────────────────────────────────────────────────────────────────┐    │
│  │                           PENDING_ACTIONS                                   │    │
│  ├─────────────────────────────────────────────────────────────────────────────┤    │
│  │  id              TEXT PRIMARY KEY                                           │    │
│  │  type            TEXT NOT NULL  -- create_payment_link | send_reminder      │    │
│  │  description     TEXT                                                       │    │
│  │  params          TEXT (JSON)                                                │    │
│  │  order_id        TEXT REFERENCES orders(id)                                 │    │
│  │  status          TEXT DEFAULT 'pending'  -- pending | approved | rejected  │    │
│  │  created_at      DATETIME DEFAULT CURRENT_TIMESTAMP                         │    │
│  │  resolved_at     DATETIME                                                   │    │
│  └─────────────────────────────────────────────────────────────────────────────┘    │
│                                                                                     │
│  ┌─────────────────────────────────────────────────────────────────────────────┐    │
│  │                           CHAT_MESSAGES                                     │    │
│  ├─────────────────────────────────────────────────────────────────────────────┤    │
│  │  id              TEXT PRIMARY KEY                                           │    │
│  │  role            TEXT NOT NULL  -- user | assistant                         │    │
│  │  content         TEXT NOT NULL                                              │    │
│  │  tool_calls      TEXT (JSON)    -- Optional: tools appelés                  │    │
│  │  created_at      DATETIME DEFAULT CURRENT_TIMESTAMP                         │    │
│  └─────────────────────────────────────────────────────────────────────────────┘    │
│                                                                                     │
└─────────────────────────────────────────────────────────────────────────────────────┘
```

### 5.2 Structures Go

```go
// internal/models/order.go
type Order struct {
    ID            string     `json:"id" db:"id"`
    ClientName    string     `json:"client_name" db:"client_name"`
    ClientEmail   string     `json:"client_email,omitempty" db:"client_email"`
    Description   string     `json:"description,omitempty" db:"description"`
    Amount        float64    `json:"amount" db:"amount"`
    Currency      string     `json:"currency" db:"currency"`
    Status        string     `json:"status" db:"status"`
    PayPalOrderID string     `json:"paypal_order_id,omitempty" db:"paypal_order_id"`
    PaymentLink   string     `json:"payment_link,omitempty" db:"payment_link"`
    CreatedAt     time.Time  `json:"created_at" db:"created_at"`
    UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`
    PaidAt        *time.Time `json:"paid_at,omitempty" db:"paid_at"`
}

// internal/models/action.go
type PendingAction struct {
    ID          string          `json:"id" db:"id"`
    Type        string          `json:"type" db:"type"`
    Description string          `json:"description" db:"description"`
    Params      json.RawMessage `json:"params" db:"params"`
    OrderID     string          `json:"order_id,omitempty" db:"order_id"`
    Status      string          `json:"status" db:"status"`
    CreatedAt   time.Time       `json:"created_at" db:"created_at"`
    ResolvedAt  *time.Time      `json:"resolved_at,omitempty" db:"resolved_at"`
}

// internal/models/message.go
type ChatMessage struct {
    ID        string          `json:"id" db:"id"`
    Role      string          `json:"role" db:"role"`
    Content   string          `json:"content" db:"content"`
    ToolCalls json.RawMessage `json:"tool_calls,omitempty" db:"tool_calls"`
    CreatedAt time.Time       `json:"created_at" db:"created_at"`
}
```

### 5.3 Types TypeScript

```typescript
// src/types/order.ts
export interface Order {
  id: string;
  client_name: string;
  client_email?: string;
  description?: string;
  amount: number;
  currency: string;
  status: 'pending' | 'paid' | 'cancelled';
  paypal_order_id?: string;
  payment_link?: string;
  created_at: string;
  updated_at: string;
  paid_at?: string;
}

// src/types/chat.ts
export interface ChatMessage {
  id: string;
  role: 'user' | 'assistant';
  content: string;
  pending_action?: PendingAction;
  created_at: string;
}

export interface PendingAction {
  id: string;
  type: 'create_payment_link' | 'send_reminder';
  description: string;
  params: Record<string, unknown>;
  status: 'pending' | 'approved' | 'rejected';
}
```

---

## 6. Spécification des APIs

### 6.1 API REST Backend

#### Base URL
```
http://localhost:8080/api (développement)
https://paypilot.example.com/api (production)
```

#### Endpoints

##### Orders

| Méthode | Endpoint | Description | Body | Réponse |
|---------|----------|-------------|------|---------|
| `GET` | `/orders` | Liste les commandes | - | `Order[]` |
| `GET` | `/orders/:id` | Détail d'une commande | - | `Order` |
| `POST` | `/orders` | Crée une commande | `CreateOrderRequest` | `Order` |
| `PATCH` | `/orders/:id` | Met à jour une commande | `UpdateOrderRequest` | `Order` |
| `DELETE` | `/orders/:id` | Supprime une commande | - | `204 No Content` |

##### Chat

| Méthode | Endpoint | Description | Body | Réponse |
|---------|----------|-------------|------|---------|
| `POST` | `/chat` | Envoie un message | `ChatRequest` | `ChatResponse` |
| `GET` | `/chat/history` | Historique des messages | - | `ChatMessage[]` |

##### Actions

| Méthode | Endpoint | Description | Body | Réponse |
|---------|----------|-------------|------|---------|
| `GET` | `/actions/pending` | Actions en attente | - | `PendingAction[]` |
| `POST` | `/actions/:id/confirm` | Confirme une action | - | `ActionResult` |
| `POST` | `/actions/:id/reject` | Rejette une action | - | `ActionResult` |

##### Analytics

| Méthode | Endpoint | Description | Query | Réponse |
|---------|----------|-------------|-------|---------|
| `GET` | `/analytics/summary` | Résumé période | `?period=month` | `Summary` |
| `GET` | `/analytics/overdue` | Commandes en retard | `?days=7` | `Order[]` |

##### PayPal

| Méthode | Endpoint | Description | Body | Réponse |
|---------|----------|-------------|------|---------|
| `POST` | `/paypal/create-order/:id` | Crée order PayPal | - | `PayPalOrder` |
| `POST` | `/paypal/capture/:paypalOrderId` | Capture paiement | - | `CaptureResult` |
| `POST` | `/paypal/webhook` | Webhook PayPal | `WebhookEvent` | `200 OK` |

### 6.2 Structures de Requêtes/Réponses

```go
// Création de commande
type CreateOrderRequest struct {
    ClientName  string  `json:"client_name" validate:"required"`
    ClientEmail string  `json:"client_email" validate:"email"`
    Description string  `json:"description"`
    Amount      float64 `json:"amount" validate:"required,gt=0"`
    Currency    string  `json:"currency" validate:"oneof=USD EUR"`
}

// Chat
type ChatRequest struct {
    Message string `json:"message" validate:"required"`
}

type ChatResponse struct {
    Message       string         `json:"message"`
    PendingAction *PendingAction `json:"pending_action,omitempty"`
}

// Analytics
type Summary struct {
    TotalRevenue     float64 `json:"total_revenue"`
    OrderCount       int     `json:"order_count"`
    PaidCount        int     `json:"paid_count"`
    PendingCount     int     `json:"pending_count"`
    PendingAmount    float64 `json:"pending_amount"`
    PreviousPeriod   float64 `json:"previous_period"`
    PercentageChange float64 `json:"percentage_change"`
}
```

---

## 7. Intégrations Externes

### 7.1 PayPal REST API v2

#### Authentification

```
POST https://api-m.sandbox.paypal.com/v1/oauth2/token

Headers:
  Authorization: Basic {base64(client_id:secret)}
  Content-Type: application/x-www-form-urlencoded

Body:
  grant_type=client_credentials
```

#### Création d'une Order PayPal

```
POST https://api-m.sandbox.paypal.com/v2/checkout/orders

Headers:
  Authorization: Bearer {access_token}
  Content-Type: application/json

Body:
{
  "intent": "CAPTURE",
  "purchase_units": [{
    "reference_id": "{order_id}",
    "amount": {
      "currency_code": "USD",
      "value": "100.00"
    },
    "description": "Site web pour Jean Dupont"
  }],
  "application_context": {
    "return_url": "https://paypilot.example.com/payment/success",
    "cancel_url": "https://paypilot.example.com/payment/cancel"
  }
}
```

#### Flux Complet

```
┌─────────────────────────────────────────────────────────────────────────────────────┐
│                              FLUX PAYPAL                                            │
├─────────────────────────────────────────────────────────────────────────────────────┤
│                                                                                     │
│  1. Création Order                                                                  │
│     Backend ──POST──▶ PayPal API ──returns──▶ { id, links: [approve_url] }          │
│                                                                                     │
│  2. Redirection Client                                                              │
│     Frontend ──redirect──▶ approve_url (PayPal Checkout)                            │
│                                                                                     │
│  3. Approbation Client                                                              │
│     Client approuve ──redirect──▶ return_url?token={order_id}                       │
│                                                                                     │
│  4. Capture Paiement                                                                │
│     Frontend ──POST──▶ Backend ──POST──▶ PayPal /capture                            │
│                                                                                     │
│  5. Webhook (optionnel mais recommandé)                                             │
│     PayPal ──POST──▶ Backend /webhook ──update──▶ Order.status = "paid"             │
│                                                                                     │
└─────────────────────────────────────────────────────────────────────────────────────┘
```

### 7.2 Claude API (Anthropic)

#### Configuration

```go
// internal/services/ai/client.go
type ClaudeClient struct {
    APIKey  string
    BaseURL string  // https://api.anthropic.com
    Model   string  // claude-sonnet-4-20250514
}
```

#### Requête avec Tools

```json
{
  "model": "claude-sonnet-4-20250514",
  "max_tokens": 1024,
  "system": "Tu es PayPilot, un assistant intelligent pour gérer les paiements...",
  "tools": [
    {
      "name": "list_orders",
      "description": "Liste les commandes avec filtres optionnels",
      "input_schema": {
        "type": "object",
        "properties": {
          "status": { "type": "string", "enum": ["pending", "paid", "cancelled"] },
          "period": { "type": "string", "enum": ["today", "this_week", "this_month"] },
          "limit": { "type": "integer", "default": 10 }
        }
      }
    }
  ],
  "messages": [
    { "role": "user", "content": "Montre-moi les commandes en attente" }
  ]
}
```

#### Flux Agent

```
┌─────────────────────────────────────────────────────────────────────────────────────┐
│                              FLUX AGENT IA                                          │
├─────────────────────────────────────────────────────────────────────────────────────┤
│                                                                                     │
│  Message User ──▶ Claude API (avec tools)                                           │
│                         │                                                           │
│                         ▼                                                           │
│               ┌─────────────────────┐                                               │
│               │  Claude décide:     │                                               │
│               │  - Réponse directe  │                                               │
│               │  - Appel de tool    │                                               │
│               └──────────┬──────────┘                                               │
│                          │                                                          │
│          ┌───────────────┼───────────────┐                                          │
│          │               │               │                                          │
│          ▼               ▼               ▼                                          │
│    Réponse directe  Tool sans      Tool avec                                        │
│    ──▶ Frontend     confirmation   confirmation                                     │
│                          │               │                                          │
│                          ▼               ▼                                          │
│                    Exécuter tool   Créer PendingAction                              │
│                          │               │                                          │
│                          ▼               ▼                                          │
│                    Résultat ──▶    Attendre confirmation                            │
│                    Claude              │                                            │
│                          │             │ [User confirme]                            │
│                          │             ▼                                            │
│                          │       Exécuter tool                                      │
│                          │             │                                            │
│                          ▼             ▼                                            │
│                    Réponse finale ◀────┘                                            │
│                          │                                                          │
│                          ▼                                                          │
│                      Frontend                                                       │
│                                                                                     │
└─────────────────────────────────────────────────────────────────────────────────────┘
```

---

## 8. Sécurité

### 8.1 Mesures Implémentées

| Aspect | Mesure |
|--------|--------|
| **CORS** | Origines autorisées configurées |
| **Validation** | Tous les inputs validés côté serveur |
| **Secrets** | Variables d'environnement, jamais en code |
| **PayPal Webhook** | Signature vérifiée |
| **SQL Injection** | Requêtes préparées (paramétrisées) |
| **XSS** | React échappe par défaut |

### 8.2 Variables d'Environnement

```env
# .env (NE PAS COMMIT)
PORT=8080
DATABASE_PATH=./paypilot.db

# PayPal Sandbox
PAYPAL_CLIENT_ID=xxx
PAYPAL_SECRET=xxx
PAYPAL_BASE_URL=https://api-m.sandbox.paypal.com

# Claude API
ANTHROPIC_API_KEY=xxx

# Frontend URL (pour CORS et redirects)
FRONTEND_URL=http://localhost:5173
```

---

## 9. Déploiement

### 9.1 Architecture de Déploiement

```
┌─────────────────────────────────────────────────────────────────────────────────────┐
│                              DÉPLOIEMENT (Render/Railway)                           │
├─────────────────────────────────────────────────────────────────────────────────────┤
│                                                                                     │
│                              ┌─────────────────┐                                    │
│                              │    Internet     │                                    │
│                              └────────┬────────┘                                    │
│                                       │                                             │
│                                       ▼                                             │
│                        ┌──────────────────────────────┐                             │
│                        │      Load Balancer/CDN       │                             │
│                        │        (Cloudflare)          │                             │
│                        └──────────────┬───────────────┘                             │
│                                       │                                             │
│              ┌────────────────────────┼────────────────────────┐                    │
│              │                        │                        │                    │
│              ▼                        ▼                        ▼                    │
│    ┌─────────────────┐    ┌─────────────────────┐    ┌─────────────────┐            │
│    │  Static Files   │    │    API Server       │    │  Webhook URL    │            │
│    │  (Vercel/       │    │    (Render/         │    │  (même API)     │            │
│    │   Netlify)      │    │     Railway)        │    │                 │            │
│    │                 │    │                     │    │                 │            │
│    │  React Build    │    │  Go Binary          │    │  /api/paypal/   │            │
│    │                 │    │  + SQLite file      │    │  webhook        │            │
│    └─────────────────┘    └─────────────────────┘    └─────────────────┘            │
│                                                                                     │
└─────────────────────────────────────────────────────────────────────────────────────┘
```

### 9.2 Options de Déploiement

| Service | Frontend | Backend | Gratuit |
|---------|----------|---------|---------|
| **Render** | Static Site | Web Service | Oui (avec limitations) |
| **Railway** | - | Container | $5 crédit gratuit |
| **Vercel** | Oui | - | Oui |
| **Fly.io** | - | Container | Oui (limité) |

**Recommandation:** Render pour tout (frontend static + backend Go)

---

## 10. Planning de Développement

### 10.1 Vue d'Ensemble (20 jours)

```
┌─────────────────────────────────────────────────────────────────────────────────────┐
│                              PLANNING 20 JOURS                                      │
├─────────────────────────────────────────────────────────────────────────────────────┤
│                                                                                     │
│  SEMAINE 1 (J1-J7)                                                                  │
│  ════════════════                                                                   │
│  ┌─────┬─────┬─────┬─────┬─────┬─────┬─────┐                                        │
│  │ J1  │ J2  │ J3  │ J4  │ J5  │ J6  │ J7  │                                        │
│  ├─────┴─────┴─────┼─────┴─────┴─────┴─────┤                                        │
│  │     SETUP       │     PAYPAL + CRUD     │                                        │
│  │  Go + React     │    Orders Backend     │                                        │
│  │  Structure      │    + Frontend         │                                        │
│  └─────────────────┴───────────────────────┘                                        │
│                                                                                     │
│  SEMAINE 2 (J8-J14)                                                                 │
│  ═════════════════                                                                  │
│  ┌─────┬─────┬─────┬─────┬─────┬─────┬─────┐                                        │
│  │ J8  │ J9  │ J10 │ J11 │ J12 │ J13 │ J14 │                                        │
│  ├─────┴─────┴─────┼─────┴─────┴─────┴─────┤                                        │
│  │   DASHBOARD     │      AGENT IA         │                                        │
│  │   Analytics     │   Claude + Tools      │                                        │
│  │   Charts        │   Chat Interface      │                                        │
│  └─────────────────┴───────────────────────┘                                        │
│                                                                                     │
│  SEMAINE 3 (J15-J20)                                                                │
│  ══════════════════                                                                 │
│  ┌─────┬─────┬─────┬─────┬─────┬─────┐                                              │
│  │ J15 │ J16 │ J17 │ J18 │ J19 │ J20 │                                              │
│  ├─────┴─────┴─────┼─────┴─────┼─────┤                                              │
│  │  HUMAN-IN-LOOP  │   POLISH  │SUBMIT                                              │
│  │  Confirmations  │   Démo    │     │                                              │
│  │  Tests          │   Vidéo   │     │                                              │
│  └─────────────────┴───────────┴─────┘                                              │
│                                                                                     │
└─────────────────────────────────────────────────────────────────────────────────────┘
```

### 10.2 Détail par Jour

| Jour | Tâches | Livrables |
|------|--------|-----------|
| **J1** | Setup Go (structure, config, DB) | Backend compilable |
| **J2** | Setup React (Vite, Tailwind, routing) | Frontend fonctionnel |
| **J3** | Modèles + migrations + repos | DB avec tables |
| **J4** | CRUD Orders (backend) | API /orders/* |
| **J5** | Client PayPal + création order | Intégration PayPal |
| **J6** | UI Orders (liste, création) | Pages orders |
| **J7** | Bouton PayPal + capture | Paiement complet |
| **J8** | Service analytics | Endpoint /analytics |
| **J9** | Dashboard UI (stats, graphiques) | Page dashboard |
| **J10** | Liste commandes en retard | Alertes overdue |
| **J11** | Client Claude + tools de base | Agent fonctionnel |
| **J12** | Endpoint /chat + exécution tools | Chat backend |
| **J13** | Interface chat (messages, input) | Chat UI |
| **J14** | Intégration chat + orders | Chat complet |
| **J15** | System PendingAction | Actions en attente |
| **J16** | UI confirmation actions | Modal confirmation |
| **J17** | Tests + bugfixes | App stable |
| **J18** | Déploiement Render | URL live |
| **J19** | Vidéo démo | Video < 3min |
| **J20** | README + soumission | Projet soumis |

---

## 11. Critères de Succès

### 11.1 MVP Minimum Viable

- [ ] Créer une commande avec montant et description
- [ ] Générer un lien de paiement PayPal (sandbox)
- [ ] Voir le statut de paiement mis à jour
- [ ] Poser une question à l'agent IA et obtenir une réponse utile
- [ ] Confirmer une action proposée par l'agent
- [ ] Dashboard avec revenus et commandes en attente

### 11.2 Nice-to-Have

- [ ] Graphique des revenus sur 6 mois
- [ ] Messages de relance générés par l'IA
- [ ] Multi-devise (USD + EUR)
- [ ] Historique complet du chat
- [ ] Export des données

---

## 12. Annexes

### 12.1 Ressources

- [PayPal Developer Documentation](https://developer.paypal.com/docs/)
- [PayPal Sandbox](https://developer.paypal.com/tools/sandbox/)
- [Anthropic Claude API](https://docs.anthropic.com/)
- [Go Chi Router](https://github.com/go-chi/chi)
- [React Documentation](https://react.dev/)
- [Tailwind CSS](https://tailwindcss.com/)
- [Zustand](https://github.com/pmndrs/zustand)

### 12.2 Glossaire

| Terme | Définition |
|-------|------------|
| **Order** | Commande créée dans PayPilot |
| **PayPal Order** | Order créée via l'API PayPal pour le paiement |
| **Tool** | Fonction que l'agent IA peut appeler |
| **Pending Action** | Action proposée par l'IA en attente de confirmation |
| **Sandbox** | Environnement de test PayPal (pas de vrais paiements) |
| **Human-in-the-loop** | Validation humaine avant exécution d'actions sensibles |

---

**Fin du Cahier des Charges**
