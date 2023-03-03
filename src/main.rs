use futures_util::stream::StreamExt;
use serde::{Deserialize, Serialize};
use std::io::{self, Write};

#[derive(Clone, Serialize, Deserialize, Debug)]
struct Message {
    role: String,
    content: String,
}

#[derive(Serialize, Deserialize, Debug)]
struct ChatCompletionRequest {
    model: String,
    messages: Vec<Message>,
    stream: bool,
}

#[derive(Serialize, Deserialize, Debug)]
struct ChatCompletionChunk {
    id: String,
    object: String,
    created: i32,
    model: String,
    choices: Vec<Choice>,
}

#[derive(Clone, Serialize, Deserialize, Debug)]
struct Choice {
    delta: Delta,
}

#[derive(Clone, Serialize, Deserialize, Debug)]
struct Delta {
    role: Option<String>,
    content: Option<String>,
}

struct Client {
    apikey: String,
    client: reqwest::Client,
}

impl Client {
    async fn complete(&self, messages: &mut Vec<Message>) {
        let completion = ChatCompletionRequest {
            model: "gpt-3.5-turbo".to_string(),
            messages: messages.to_vec(),
            stream: true,
        };

        let req = self
            .client
            .post("https://api.openai.com/v1/chat/completions")
            .header("Content-Type", "application/json")
            .header("Authorization", format!("Bearer {}", self.apikey,));
        let body = serde_json::to_string(&completion).unwrap();
        let req = req.body(body);

        let mut stream = req.send().await.unwrap().bytes_stream();
        let mut completion_output = String::new();

        while let Some(item) = stream.next().await {
            let item = item.unwrap();
            let vec = &item.to_vec();
            let text = std::str::from_utf8(vec).unwrap();

            for data in text.split("data: ") {
                if data.contains("[DONE]") {
                    break;
                }
                if data.len() < 6 {
                    continue;
                }

                let data = data.trim();
                let data = match serde_json::from_str::<ChatCompletionChunk>(data) {
                    Ok(data) => data,
                    Err(e) => {
                        println!("error on deserializing chunk: {:#?}", e);
                        continue;
                    }
                };

                let choice = data.choices[0].clone();
                let content = &choice.delta.content.unwrap_or("".to_string());
                print!("{}", content);
                completion_output.push_str(content);
                io::stdout().flush().unwrap();
            }
        }

        messages.push(Message {
            role: "system".to_string(),
            content: completion_output,
        });
        println!("");
    }
}

#[tokio::main]
async fn main() {
    let client = Client {
        apikey: "sk-xxxxxxxxxxxxxxxxxxxxxxxx".to_string(),
        client: reqwest::Client::new(),
    };

    let mut messages: Vec<Message> = vec![];

    loop {
        print!("> ");
        io::stdout().flush().unwrap();
        let mut input = String::new();
        io::stdin().read_line(&mut input).unwrap();

        messages.push(Message {
            role: "user".to_string(),
            content: input.clone(),
        });

        client.complete(&mut messages).await;
    }
}
