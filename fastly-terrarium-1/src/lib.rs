#[macro_use]
extern crate http_guest;

use http_guest::{Request, Response};

pub fn user_entrypoint(_req: &Request<Vec<u8>>) -> Response<Vec<u8>> {
    Response::builder()
        .status(200)
        .body(_req.headers().get("fastly-client-ip").unwrap().as_bytes().to_owned())
        .unwrap()
}

guest_app!(user_entrypoint);
