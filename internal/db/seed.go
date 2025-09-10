package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"

	"github.com/orangeMangoDimz/go-social/internal/store"
)

var usernames = []string{
	"CrimsonGhost",
	"SteelNinja",
	"QuantumLeap",
	"ShadowWalker",
	"SolarFlare",
	"CyberKnight",
	"MysticEcho",
	"ZeroGravity",
	"IronWarlock",
	"AzureDream",
	"NeonSpecter",
	"RoguePhoenix",
	"DigitalNomad",
	"AtomicPulse",
	"VelvetThunder",
	"CosmicRider",
	"FrostWraith",
	"CrimsonFury",
	"SiliconSage",
	"EchoWhisperer",
	"PhantomByte",
	"StarGazerX",
	"GravityShift",
	"TurboPenguin",
	"BlazingComet",
	"NightHawk22",
	"QuantumJolt",
	"MysticRaven7",
	"SolarPioneer",
	"VoidRunner",
	"CodeBreaker01",
	"IroncladDev",
	"PixelProwler",
	"ShadowMatrix",
	"GalacticDrifter",
	"NovaSurge",
	"StaticShock",
	"CrimsonCipher",
	"AzureStriker",
	"TerraByte",
	"EchoRider",
	"QuantumPulse",
	"DarkMatter",
	"SolarWind",
	"ByteSlinger",
	"GhostProtocol",
	"ZenMaster",
	"ApexPredator",
	"RogueLogic",
	"CobaltCrusader",
}

var titles = []string{
	"Mastering Concurrency in Go: A Practical Guide",
	"Python's GIL Explained: Unlocking Performance",
	"Building Scalable Microservices with Go and gRPC",
	"Optimizing Database Queries in Python Applications",
	"A Deep Dive into Go's Garbage Collector",
	"Asynchronous Python with AsyncIO: A Complete Guide",
	"Deploying Go Applications with Docker and Kubernetes",
	"REST vs. GraphQL: Choosing the Right API Architecture",
	"Effective Error Handling Patterns in Go",
	"Dependency Injection in Python with FastAPI",
	"Securing Your Go Backend: Best Practices",
	"From Monolith to Microservices: A Backend Engineer's Journey",
	"Performance Profiling in Go: Finding Your Bottlenecks",
	"The Art of Writing Clean, Idiomatic Python",
	"Using Channels in Go for Powerful Concurrency",
	"Caching Strategies for High-Performance Backends",
	"Building a Real-Time Chat Application with Go and WebSockets",
	"Python Type Hinting: A Game-Changer for Code Quality",
	"Architecting Resilient Systems in a Distributed World",
	"Test-Driven Development (TDD) in Go: A Step-by-Step Guide",
	"Working with PostgreSQL Efficiently in Python",
	"Go Modules: Managing Dependencies Like a Pro",
	"A Practical Introduction to CI/CD for Backend Engineers",
	"How We Reduced API Latency by 50% with Go",
	"Understanding Interfaces in Go: A Simple Analogy",
}

var contents = []string{
	"Go's concurrency model, built on goroutines and channels, is a game-changer for building high-performance applications. This post is a practical deep-dive into creating lightweight goroutines with the `go` keyword and using channels for safe communication and synchronization between them. We'll explore core patterns like worker pools and the select statement for handling multiple channel operations. Learn how to write clean, concurrent code while avoiding common pitfalls such as race conditions and deadlocks. By the end, you'll have a solid foundation for leveraging Go's powerful concurrency primitives to build faster, more responsive backend systems.",
	"The Global Interpreter Lock (GIL) is a famous, and often misunderstood, feature of CPython. It's a mutex that ensures only one thread executes Python bytecode at a time, which simplifies memory management but limits the performance of multi-threaded, CPU-bound programs. This article clarifies what the GIL is, why it exists, and its real-world impact. We'll discuss the difference between CPU-bound and I/O-bound tasks and show why the GIL isn't a problem for the latter. We will also explore effective strategies to mitigate its effects, including using the `multiprocessing` module for true parallelism and considering alternative Python interpreters.",
	"Microservice architectures demand efficient and robust inter-service communication. This is where gRPC shines. Built on HTTP/2 and using Protocol Buffers for schema definition, gRPC offers significant performance benefits over traditional REST APIs. This guide walks you through building a microservice in Go using gRPC from scratch. You'll learn how to define your service and message types in a `.proto` file, generate Go code, implement the server logic, and create a client to consume it. We'll also cover advanced features like streaming, error handling, and metadata, giving you the tools to build truly scalable backend systems.",
	"Slow database queries are a primary cause of poor application performance. This guide provides actionable techniques for optimizing database interactions in your Python applications. We'll cover everything from fundamental indexing strategies (B-Tree, Hash) to writing efficient SQL queries and avoiding common ORM pitfalls like the N+1 problem. Learn how to use tools like `EXPLAIN ANALYZE` to diagnose slow queries and understand their execution plans. We'll also discuss the importance of connection pooling and caching strategies (e.g., using Redis) to reduce database load and significantly improve your application's response times.",
	"Go's automatic memory management simplifies development, but understanding how the garbage collector (GC) works is crucial for building high-performance, low-latency applications. This article provides a deep dive into Go's concurrent, tri-color mark-and-sweep garbage collector. We'll break down the GC cycle, explaining the mark, sweep, and pause phases. You'll learn how the GC is designed to minimize 'stop-the-world' pauses and how you can influence its behavior through tuning parameters like `GOGC`. Understanding these mechanics will help you write more memory-efficient code and diagnose potential performance issues related to memory allocation.",
	"Asynchronous programming is essential for building I/O-bound applications that can handle many connections simultaneously. Python's `asyncio` library provides a powerful framework for writing single-threaded concurrent code using the `async/await` syntax. This guide will take you from the basics of coroutines and the event loop to advanced patterns for managing thousands of concurrent tasks. We'll cover how to work with asynchronous libraries for databases and HTTP requests, and how to properly handle exceptions and task cancellation in an async environment. Unlock the power of non-blocking I/O in your Python backends.",
	"Containerization has revolutionized how we build, ship, and run applications. This comprehensive tutorial provides a step-by-step guide to deploying your Go applications using Docker and Kubernetes. First, we'll cover how to write an efficient, multi-stage Dockerfile to create a small and secure image for your Go binary. Then, we'll dive into Kubernetes, showing you how to write YAML manifests for Deployments, Services, and Ingress to run your application at scale. Learn how to manage configuration with ConfigMaps and Secrets, and how to set up health checks for a resilient, self-healing system.",
	"Choosing the right API architecture is a critical decision in modern backend development. This post offers a detailed comparison between REST (Representational State Transfer) and GraphQL. We'll break down the core principles of each, comparing how they handle data fetching, versioning, and schema definition. Explore REST's simplicity and widespread adoption versus GraphQL's power to eliminate over-fetching and under-fetching with its typed query language. We provide clear use-cases and examples in both Go and Python to help you decide which technology is the best fit for your next project, considering factors like client requirements, performance, and developer experience.",
	"Robust error handling is a hallmark of professional Go code. Unlike languages with exceptions, Go uses an explicit, value-based error handling model. This article explores best practices and patterns for managing errors effectively in your Go applications. We'll go beyond the simple `if err != nil` check and discuss techniques for adding context to errors, defining custom error types, and handling errors gracefully at API boundaries. Learn how to use the `errors` package, including features like `errors.Is()` and `errors.As()`, to create maintainable and debuggable code that anticipates and communicates failure clearly.",
	"Dependency Injection (DI) is a powerful design pattern for building loosely coupled, maintainable, and testable applications. FastAPI, a modern Python web framework, has first-class support for DI that makes it incredibly easy to implement. This guide explains the core concepts of DI and demonstrates how to leverage FastAPI's `Depends` system. Learn how to inject dependencies like database sessions, configuration objects, and service classes directly into your path operation functions. We'll show you how this not only cleans up your code but also simplifies unit testing by allowing you to easily swap out dependencies with mocks.",
	"Security is not an afterthought; it's a critical component of backend development. This guide outlines essential security best practices for your Go applications. We'll cover topics such as validating user input to prevent injection attacks (SQLi, XSS), securely managing secrets and credentials, and implementing proper authentication and authorization flows (e.g., using JWTs). Learn how to protect your API endpoints with rate limiting and middleware, and why you should always use HTTPS. We'll also discuss common vulnerabilities and provide code snippets and library recommendations to help you build a more secure, hardened Go backend from day one.",
	"Migrating from a monolithic architecture to microservices is a significant undertaking that can yield incredible benefits in scalability and team autonomy. This article shares a backend engineer's journey through this complex process. We'll discuss the key drivers for making the change and the strategies used for decomposition, like the Strangler Fig pattern. Learn about the challenges encountered along the way, including data consistency across services, inter-service communication, and the need for robust observability (logging, metrics, tracing). This is a practical story filled with lessons learned and advice for any team considering a similar architectural evolution.",
	"When your Go application is slow, guessing where the bottleneck is won't work. You need data. Go comes with a powerful suite of profiling tools, including `pprof`, that can help you pinpoint performance issues. This hands-on guide will show you how to integrate `pprof` into your web server to collect and analyze CPU and memory profiles. Learn how to generate and interpret flame graphs to visualize where your program is spending its time. We'll walk through a practical example of identifying a CPU-intensive function and optimizing it, demonstrating how profiling is an essential skill for any serious Go developer.",
	"Writing code that works is one thing; writing code that is clean, readable, and maintainable is another. This article explores the art of writing idiomatic Python. We'll move beyond the basic syntax and delve into the principles outlined in PEP 8 and the Zen of Python. Learn how to use list comprehensions effectively, understand the difference between `==` and `is`, and leverage features like generators for memory-efficient iteration. We'll cover best practices for structuring your modules and packages, writing clear docstrings, and using context managers (`with` statement) to create robust and elegant Python code.",
	"Channels are the heart of Go's concurrency model, providing a typed conduit through which you can send and receive values with the channel operator, `<-`. This enables safe and easy communication between goroutines. This guide explores the fundamental patterns for using channels effectively. We'll cover the difference between unbuffered and buffered channels, how to use range over a channel to process a stream of data, and how the `select` statement allows a goroutine to wait on multiple communication operations. Learn patterns like fan-in/fan-out to build complex data processing pipelines with clean, readable concurrent code.",
	"In high-traffic systems, your database is often the bottleneck. A well-implemented caching layer is one of the most effective ways to improve performance and reduce database load. This article explores various caching strategies for backend systems. We'll discuss in-memory caching for speed, as well as distributed caching with tools like Redis for scalability and persistence. Learn about different caching patterns like cache-aside, read-through, and write-through, and understand their respective trade-offs. We'll also cover crucial topics like cache invalidation, setting appropriate TTLs (Time-To-Live), and avoiding common caching pitfalls.",
	"WebSockets provide a full-duplex communication channel over a single TCP connection, making them perfect for building interactive, real-time applications. This tutorial will guide you through the process of building a live chat application from scratch using Go for the backend. We'll use the popular `gorilla/websocket` library to handle the WebSocket connections. You'll learn how to manage multiple client connections, broadcast messages to all connected clients, and handle connection lifecycle events gracefully. This is a hands-on project that will give you a solid understanding of how to leverage WebSockets in your Go applications.",
	"Python is a dynamically typed language, but the introduction of type hints (PEP 484) has been a revolutionary addition for building large, complex applications. This article explains why you should be using type hints in your Python projects. We'll show how they dramatically improve code readability and maintainability, allowing static analysis tools like `mypy` to catch bugs before you even run the code. Learn how type hints enable better editor support with features like autocompletion and refactoring. We'll cover the basic syntax for typing variables, function arguments, and return values, demonstrating how it makes your code more robust and easier to reason about.",
	"In a world of microservices and cloud infrastructure, failure is inevitable. Resilient systems are designed with this in mind, able to withstand and recover from failures gracefully. This article discusses key architectural patterns for building resilience. We'll explore the Circuit Breaker pattern to prevent cascading failures, the importance of timeouts and retries for transient network issues, and the role of health checks in automated recovery. Learn how to design for statelessness where possible and ensure data consistency in a distributed environment. These principles are essential for building backend systems with high availability.",
	"Test-Driven Development (TDD) is a software development process that encourages simple design and inspires confidence. It revolves around a short, repetitive cycle: write a failing test, write the minimum code to pass the test, and then refactor. This guide provides a practical walkthrough of applying TDD in a Go environment. Using Go's built-in `testing` package, we will build a small application piece by piece, starting with the tests. You'll see how this process leads to a more robust and modular design, improves test coverage, and makes your code easier to maintain and refactor in the long run.",
	"PostgreSQL is a powerful, open-source object-relational database system with a strong reputation for reliability and data integrity. This guide focuses on best practices for working with PostgreSQL from a Python backend. We'll compare popular libraries like `psycopg2` and `asyncpg`, and demonstrate how to use them effectively with and without an ORM like SQLAlchemy. Learn how to manage connection pools for better performance, execute transactions safely to ensure data consistency, and leverage some of PostgreSQL's advanced features like JSONB data types and full-text search directly from your Python code.",
	"Dependency management is a critical part of any modern software project. Since Go 1.11, modules have been the official and built-in way to manage dependencies in Go. This article will get you up to speed on everything you need to know to manage your project's dependencies like a pro. We'll cover how to initialize a new module, add, update, and remove dependencies using `go get`, and understand the roles of the `go.mod` and `go.sum` files. Learn about semantic versioning and how Go uses it to ensure your builds are reproducible and stable. Say goodbye to the old GOPATH way of working and embrace the power of Go modules.",
	"Continuous Integration (CI) and Continuous Deployment (CD) are practices that automate the software release process, enabling teams to deliver code more frequently and reliably. This article provides a practical introduction to CI/CD from a backend engineer's perspective. We'll break down the stages of a typical pipeline: building the code, running tests (unit, integration), performing security scans, and deploying to various environments (staging, production). Using a popular tool like GitHub Actions, we'll walk through creating a simple CI/CD pipeline for a Go or Python application, demonstrating how automation can drastically improve your development workflow.",
	"API latency is a critical user-facing metric. A slow API leads to a poor user experience. This case study details how our team identified and resolved performance bottlenecks in a high-traffic API, resulting in a 50% reduction in average response time. We'll walk through the process, starting with establishing a baseline using monitoring and observability tools. Then we'll cover the profiling techniques used to pinpoint the exact areas of inefficient code in our Go backend. Finally, we'll discuss the specific optimizations that were implemented, from query tuning and caching strategies to concurrency pattern adjustments. Learn from our real-world experience.",
	"Interfaces are one of Go's most powerful and elegant features. They provide a way to specify behavior without being tied to a specific implementation, enabling clean, decoupled architectures. However, they can be confusing for newcomers. This post explains interfaces using a simple, real-world analogy. We'll show how interfaces in Go are satisfied implicitly, unlike in other languages. You'll learn the common use cases for interfaces, such as creating mocks for testing and writing functions that can operate on different types. We'll also cover best practices, like the idea that smaller interfaces are better ('accept interfaces, return structs').",
}

var tags = []string{
	"golang",
	"python",
	"backend",
	"microservices",
	"api-design",
	"grpc",
	"performance",
	"concurrency",
	"database",
	"postgresql",
	"docker",
	"kubernetes",
	"devops",
	"system-design",
	"architecture",
	"testing",
	"security",
	"optimization",
	"scalability",
	"clean-code",
}

func Seed(store store.Storage, db *sql.DB) {
	ctx := context.Background()

	users := generateUsers(100)
	tx, _ := db.BeginTx(ctx, nil)

	for _, user := range users {
		if err := store.Users.Create(ctx, tx, user); err != nil {
			_ = tx.Rollback()
			log.Println("Error creating users:", err)
			return
		}
	}

	tx.Commit()

	posts := generatePosts(200, users)
	for _, post := range posts {
		if err := store.Posts.Create(ctx, post); err != nil {
			log.Println("Error creating post:", err)
			return
		}
	}

	comments := generateComments(500, users, posts)
	for _, comment := range comments {
		if err := store.Comments.Create(ctx, comment); err != nil {
			log.Println("Error creating comment:", err)
			return
		}
	}

	log.Println("Seeding complete!")
}

func generateUsers(num int) []*store.User {
	users := make([]*store.User, num)

	for i := 0; i < num; i++ {
		users[i] = &store.User{
			Username: usernames[i%len(usernames)] + fmt.Sprintf("%d", i+1),
			Email:    usernames[i%len(usernames)] + fmt.Sprintf("%d", i+1) + "@example.com",
			Role: store.Role{
				Name: "user",
			},
		}
	}
	return users
}

func generatePosts(num int, users []*store.User) []*store.Post {
	posts := make([]*store.Post, num)
	for i := 0; i < num; i++ {
		user := users[rand.Intn(len(users))]
		posts[i] = &store.Post{
			UserId:  user.ID,
			Title:   titles[rand.Intn(len(titles))],
			Content: contents[rand.Intn(len(contents))],
			Tags: []string{
				tags[rand.Intn(len(tags))],
				tags[rand.Intn(len(tags))],
			},
		}
	}

	return posts
}

func generateComments(num int, users []*store.User, posts []*store.Post) []*store.Comment {
	comments := make([]*store.Comment, num)
	for i := 0; i < num; i++ {
		comments[i] = &store.Comment{
			PostID:  posts[rand.Intn(len(posts))].ID,
			UserID:  users[rand.Intn(len(users))].ID,
			Content: contents[rand.Intn(len(contents))],
		}
	}

	return comments
}
