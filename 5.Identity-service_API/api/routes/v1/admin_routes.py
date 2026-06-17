from fastapi import APIRouter

router = APIRouter()

@router.get("/")
async def index():
    return {"message": "Hello Admin!"}

@router.get("/users")
async def users():

    return {"message": "Hello Admin!"}