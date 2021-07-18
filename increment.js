exports = async function(changeEvent) {
    var docId = changeEvent.fullDocument._id;
    
    const countercollection = context.services.get("Cloud0").db(changeEvent.ns.db).collection("counters");
    const postsscollection = context.services.get("Cloud0").db(changeEvent.ns.db).collection(changeEvent.ns.coll);
    
    var counter = await countercollection.findOneAndUpdate({_id: changeEvent.ns },{ $inc: { seq_value: 1 }}, { returnNewDocument: true, upsert : true});
    var updateRes = await postscollection.updateOne({_id : docId},{ $set : {id : counter.seq_value}});
    
    console.log(`Updated ${JSON.stringify(changeEvent.ns)} with counter ${counter.seq_value} result : ${JSON.stringify(updateRes)}`);
    };