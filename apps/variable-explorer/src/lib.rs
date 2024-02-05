use spin_sdk::http::{IntoResponse, Request, Response};
use spin_sdk::{http_component, variables};

/// A simple Spin HTTP component.
#[http_component]
fn handle_variable_explorer(_req: Request) -> anyhow::Result<impl IntoResponse> {
    let log_level = variables::get("log_level")?;
    let platform_name = variables::get("platform_name")?;
    let db_password = variables::get("db_password")?;

    println!("# Log Level: {}", log_level);
    println!("# Platform name: {}", platform_name);
    println!("# DB Password: {}", db_password);

    Ok(Response::builder()
        .status(200)
        .header("content-type", "text/plain")
        .body(format!("Hell from {}", platform_name))
        .build())
}
